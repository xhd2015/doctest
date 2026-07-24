# Metrics Foundation — Project ID, Run Paths, JSONL Writer

## Version
0.0.2

Classic-TDD specification for the greenfield package
`github.com/xhd2015/doctest/libdoc/metrics` (P1 only).

These tests call package APIs directly. They do **not** wire `doctest test`,
CLI metrics commands, WARNING banners, or the `review-perf` skill.

Until the implementer creates `libdoc/metrics`, this tree is expected to **RED**
(compile/import failure or failing assertions).

# DSN (Domain Specific Notion)

### Participants

- **Caller** — a future `doctest test` process (or these unit-style leaves) that
  records run and leaf timing without coupling to the runner internals.
- **Project identity** — a stable slug for the repository under test, derived
  from `git remote get-url origin` when available, else a hash of the absolute
  module/root path.
- **Cache layout** — durable files under
  `$CACHE/doctest/metrics/<project_id>/runs/`.
- **Run file** — one append-only JSONL file owned by a single process for one
  test run; filename encodes UTC time, a two-digit disambiguator, and a unique
  suffix.
- **Buffered writer** — mutex-protected memory buffer that appends JSON lines
  (`\n`-terminated). Flushes when the buffer reaches **128 KiB** or when the
  writer is closed. No fsync requirement.
- **Event stream** — ordered JSON objects with `schema_version: 1` and a `type`
  discriminator: `run_start`, `leaf_start`, `leaf_end`, `run_end`.

### Behaviors

- **Resolve project id** — normalize origin URL → host/path without scheme,
  userinfo, or trailing `.git`; slugify `/` → `_`. When origin is missing/empty,
  use `nogit_<sha256(abs_root)[:12]>`.
- **Allocate run path** — build
  `.../runs/YYYY-MM-DD-HH-MM-SS-NN-<suffix>.jsonl` in UTC; exclusive-create
  retries `NN` (00–99) or varies suffix so two opens never clobber.
- **Encode events** — marshal typed events to one JSON object per line.
  `leaf_end` with `result=skip` may omit a matching `leaf_start` and may omit
  `ts_start`. `run_start` may carry git branch/commit but **must not** include
  a dirty/worktree flag.
- **Buffer and flush** — small writes stay in memory until `Close`; cumulative
  buffer ≥ 128 KiB forces a mid-run flush so the file is non-empty before close.
- **Partial runs** — a file closed (or left) without `run_end` remains readable
  as a partial JSONL prefix.

### Pipeline sketch

```
origin or abs-root
  -> project_id slug
  -> $CACHE/doctest/metrics/<project_id>/runs/<utc>-NN-suffix.jsonl
  -> Writer.Write(event*) [buffer 128KiB]
  -> Close flush
  -> JSONL lines: run_start / leaf_* / run_end
```

## Decision Tree

```
tests/metrics-foundation/
├── project-id/                         [identity: origin vs fallback]
│   ├── https-origin/                   HTTPS GitHub URL → github.com_xhd2015_doctest
│   ├── ssh-origin/                     git@host:path → same slug
│   └── fallback-no-origin/             empty origin → nogit_<12 hex>
├── run-path/                           [path naming under cache]
│   ├── utc-filename-pattern/           UTC YYYY-MM-DD-HH-MM-SS-NN-suffix.jsonl
│   └── exclusive-create-disambiguates/ two creates same second → distinct files
├── events/                             [schema_version 1 encode/decode]
│   ├── full-lifecycle-roundtrip/       run_start → leaf_start → leaf_end(pass) → leaf_end(skip) → run_end
│   ├── skip-as-leaf-end-only/          skip recorded as leaf_end only
│   └── run-start-omits-dirty/          run_start has no dirty/git-dirty field
└── writer/                             [128KiB buffer + close]
    ├── flush-on-close-small/           small events empty on disk until Close
    ├── flush-at-128kib/                ≥128KiB forces flush before Close
    └── partial-without-run-end/        readable partial without run_end
```

## Test Index

| Leaf | Scenario |
|------|----------|
| `project-id/https-origin` | HTTPS origin slugifies to `github.com_xhd2015_doctest` |
| `project-id/ssh-origin` | SSH scp-like origin yields the same slug |
| `project-id/fallback-no-origin` | No origin → `nogit_` + 12 hex of sha256(abs root) |
| `run-path/utc-filename-pattern` | Path under cache matches UTC run filename pattern |
| `run-path/exclusive-create-disambiguates` | Two exclusive creates do not share a path |
| `events/full-lifecycle-roundtrip` | Write five event types; re-read valid JSONL |
| `events/skip-as-leaf-end-only` | Skip is a `leaf_end` with `result=skip` (no start required) |
| `events/run-start-omits-dirty` | Decoded `run_start` has no dirty-worktree field |
| `writer/flush-on-close-small` | Sub-threshold buffer hits disk only on Close |
| `writer/flush-at-128kib` | ≥128 KiB buffered payload flushes mid-run |
| `writer/partial-without-run-end` | File without `run_end` still decodes prior lines |

## How to Run

```sh
doctest vet ./tests/metrics-foundation/
doctest test ./tests/metrics-foundation/          # RED until libdoc/metrics exists
doctest test ./tests/metrics-foundation/project-id/...
doctest test ./tests/metrics-foundation/run-path/...
doctest test ./tests/metrics-foundation/events/...
doctest test ./tests/metrics-foundation/writer/...
```

```go
import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/metrics"
)

// Request selects one metrics API surface. Leaves set Op and related fields.
type Request struct {
	Op string // project_id_from_origin | project_id_fallback | run_file_path | create_run_files | write_sequence

	// project id
	Origin  string
	AbsRoot string

	// paths
	CacheDir  string
	ProjectID string
	At        time.Time
	NN        int
	Suffix    string
	CreateCount int // exclusive create attempts (same At)

	// writer sequence
	Events              []map[string]any // encoded as schema_version 1 events via Writer
	PadBytes            int              // if >0, append one large raw JSON object (~PadBytes)
	InspectBeforeClose  bool
	CloseWriter         bool // default true when Op=write_sequence unless set false via leave-open
	LeaveOpen           bool // if true, do not Close (and do not require mid inspect after close)
}

type Response struct {
	ProjectID string

	Path  string
	Paths []string

	SizeBeforeClose int64
	SizeAfterClose  int64
	FileData        []byte
	Lines           []string
	Decoded         []map[string]any

	Err string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}

	switch req.Op {
	case "project_id_from_origin":
		resp.ProjectID = metrics.ProjectIDFromOrigin(req.Origin)
		return resp, nil

	case "project_id_fallback":
		resp.ProjectID = metrics.ProjectIDFallback(req.AbsRoot)
		return resp, nil

	case "run_file_path":
		resp.Path = metrics.RunFilePath(req.CacheDir, req.ProjectID, req.At, req.NN, req.Suffix)
		return resp, nil

	case "create_run_files":
		n := req.CreateCount
		if n <= 0 {
			n = 1
		}
		at := req.At
		if at.IsZero() {
			at = time.Now().UTC()
		}
		for i := 0; i < n; i++ {
			p, err := metrics.CreateRunFile(req.CacheDir, req.ProjectID, at, req.Suffix)
			if err != nil {
				resp.Err = err.Error()
				return resp, err
			}
			resp.Paths = append(resp.Paths, p)
		}
		if len(resp.Paths) > 0 {
			resp.Path = resp.Paths[0]
		}
		return resp, nil

	case "write_sequence":
		if req.CacheDir == "" {
			req.CacheDir = t.TempDir()
		}
		if req.ProjectID == "" {
			req.ProjectID = "test_project"
		}
		at := req.At
		if at.IsZero() {
			at = time.Now().UTC()
		}
		path, err := metrics.CreateRunFile(req.CacheDir, req.ProjectID, at, req.Suffix)
		if err != nil {
			// Fall back to explicit path + open if CreateRunFile is path-only in some implementations
			path = metrics.RunFilePath(req.CacheDir, req.ProjectID, at, 0, orSuffix(req.Suffix))
			if mkErr := os.MkdirAll(filepath.Dir(path), 0o755); mkErr != nil {
				resp.Err = mkErr.Error()
				return resp, mkErr
			}
		}
		resp.Path = path

		w, err := metrics.OpenWriter(path)
		if err != nil {
			resp.Err = err.Error()
			return resp, err
		}

		for _, ev := range req.Events {
			if ev == nil {
				continue
			}
			if _, ok := ev["schema_version"]; !ok {
				ev["schema_version"] = metrics.SchemaVersion
			}
			if err := w.Write(ev); err != nil {
				_ = w.Close()
				resp.Err = err.Error()
				return resp, err
			}
		}
		if req.PadBytes > 0 {
			// One JSON object large enough to push the buffer over the flush threshold.
			pad := strings.Repeat("x", req.PadBytes)
			padEv := map[string]any{
				"schema_version": metrics.SchemaVersion,
				"type":           "pad",
				"blob":           pad,
			}
			if err := w.Write(padEv); err != nil {
				_ = w.Close()
				resp.Err = err.Error()
				return resp, err
			}
		}

		if req.InspectBeforeClose {
			if st, stErr := os.Stat(path); stErr == nil {
				resp.SizeBeforeClose = st.Size()
			} else if !os.IsNotExist(stErr) {
				_ = w.Close()
				resp.Err = stErr.Error()
				return resp, stErr
			}
			// missing file ⇒ size 0
		}

		if !req.LeaveOpen {
			if err := w.Close(); err != nil {
				resp.Err = err.Error()
				return resp, err
			}
			if st, stErr := os.Stat(path); stErr != nil {
				resp.Err = stErr.Error()
				return resp, stErr
			} else {
				resp.SizeAfterClose = st.Size()
			}
			data, rErr := os.ReadFile(path)
			if rErr != nil {
				resp.Err = rErr.Error()
				return resp, rErr
			}
			resp.FileData = data
			resp.Lines, resp.Decoded = splitJSONL(data)
		}
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown Op %q", req.Op)
	}
}

func orSuffix(s string) string {
	if s != "" {
		return s
	}
	return "deadbeef"
}

func splitJSONL(data []byte) (lines []string, decoded []map[string]any) {
	sc := bufio.NewScanner(bytes.NewReader(data))
	// allow long pad lines
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 4*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		lines = append(lines, line)
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err == nil {
			decoded = append(decoded, m)
		} else {
			decoded = append(decoded, map[string]any{"_parse_error": err.Error(), "_raw": line})
		}
	}
	return lines, decoded
}
```
