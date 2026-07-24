# Metrics Record — doctest test JSONL + default-suite WARNING

## Version
0.0.2

**Layer model (coverage backfill):** this tree is **L2 doctest in-process**.
Leaves call pure `libdoc/metrics` helpers, `runner.ParseTestOptions`, and
`runner.RunTest` in the same process as the harness — **no product binary
spawn**. Recording leaves may still take multi-second wall time (nested
prepare / go test via `RunTest`) and carry `label: heavy` accordingly.

Coverage for metrics recording and the default-suite slow WARNING. Depends on
package `github.com/xhd2015/doctest/libdoc/metrics` (writer, paths, project id)
and runner/core options (`MetricsOn`, `MetricsRoot`).

Out of scope: `doctest metrics` analyze (`tests/metrics-cli`), review-perf skill
docs, env silence for WARNING, fail-on-slow-suite.

# DSN (Domain Specific Notion)

### Participants

- **Caller** — package entry `runner.RunTest` (or pure helpers) that owns one
  suite wall clock.
- **Warn helpers** — pure functions deciding whether to emit the default-suite
  slow WARNING and formatting its fixed message (no real 3-minute sleep in tests).
- **Metrics options** — `MetricsOn` (CLI `--metrics-on` / `ParseTestOptions`) and
  injectable `MetricsRoot` (test cache root; production uses the user cache layout).
- **Run recorder** — opens one exclusive JSONL run file under
  `$MetricsRoot/doctest/metrics/<project_id>/runs/`, writes `run_start` /
  leaf events / `run_end` via the buffered writer, closes on suite end.
- **Default suite** — discovery with `!LabelAll && len(LabelExprs)==0`.
- **Stderr** — receives non-fatal `WARNING:` after the result summary when the
  warn predicate fires; process exit is never driven solely by this warning.

### Behaviors

- **Should warn** — true only when default suite, `total > 0`, and
  `elapsed > DefaultSuiteWarnThreshold` (3 minutes). False when labeled /
  `--label-all`, when elapsed ≤ threshold, or when total is 0.
- **Format warning** — fixed prose mentioning `WARNING:`, default suite speed
  (3 minutes), `skill:doctest-review-perf`, and `doctest skill review-perf --show`.
- **Opt-in** — metrics are off by default; `--metrics-on` / `MetricsOn` enables
  creating a run file under MetricsRoot.
- **Record** — with metrics on: write `run_start` (project_id, cwd, argv/flags,
  git branch/commit without dirty, session_id, `mode.default_suite`,
  `schema_version` 1), leaf_start/leaf_end for executed leaves when practical,
  then `run_end` (wall_ns, passed, total, skipped, exit_ok, warnings including
  `default_suite_slow` when the warn predicate fired).
- **Emit WARNING** — after suite finish (after result summary preferred); yellow
  when color is enabled (color itself is not asserted in this tree).

### Pipeline sketch

```
ParseTestOptions / core.Options
  -> MetricsOn, MetricsRoot, labels
  -> runner.RunTest(dir, opts)   # in-process
  -> if MetricsOn: open run JSONL under MetricsRoot
  -> run_start
  -> execute leaves -> leaf_start / leaf_end*
  -> run_end (+ warnings)
  -> PrintResultSummary
  -> if ShouldWarnDefaultSuiteSlow: FormatDefaultSuiteSlowWarning -> stderr
```

## Decision Tree

```
tests/metrics-record/                          [L2 in-process]
├── warn-predicate/                            [pure ShouldWarnDefaultSuiteSlow]
│   ├── fires-default-over-threshold/          default + total>0 + elapsed>3m → true
│   ├── no-fire-elapsed-at-or-under/           default + elapsed≤3m → false
│   ├── no-fire-label-all/                     LabelAll + slow → false
│   ├── no-fire-label-exprs/                   LabelExprs set + slow → false
│   └── no-fire-total-zero/                    total==0 + slow → false
├── warning-message/                           [FormatDefaultSuiteSlowWarning]
│   └── required-phrases/                      WARNING:, skill, review-perf --show, 3 minutes
├── flags/                                     [ParseTestOptions]
│   ├── default-metrics-on/                    omit --metrics-on → MetricsOn=false
│   └── no-metrics-sets-opt-out/               --metrics-on → MetricsOn=true
└── recording/                                 [runner.RunTest + MetricsRoot]
    ├── no-metrics-writes-nothing/             MetricsOn=false → no new *.jsonl
    ├── enabled-writes-run-start-end/          MetricsOn=true → run_start + run_end  [heavy]
    ├── enabled-writes-leaf-events/            MetricsOn=true + 1 leaf → leaf_*      [heavy]
    └── phases/                                phase events via RunTest
        ├── emits-tree-phases/
        └── leaf-not-full-tree-wall/
```

## Test Index

| Leaf | Layer | Expected |
|------|--------|----------|
| `warn-predicate/fires-default-over-threshold` | L2 pure | `ShouldWarn...` true |
| `warn-predicate/no-fire-elapsed-at-or-under` | L2 pure | false at 3m and under 3m |
| `warn-predicate/no-fire-label-all` | L2 pure | false when not default (label-all) |
| `warn-predicate/no-fire-label-exprs` | L2 pure | false when LabelExprs non-empty |
| `warn-predicate/no-fire-total-zero` | L2 pure | false when total=0 |
| `warning-message/required-phrases` | L2 pure | fixed message substrings |
| `flags/default-metrics-on` | L2 parse | `MetricsOn == false` by default |
| `flags/no-metrics-sets-opt-out` | L2 parse | `--metrics-on` → `MetricsOn == true` |
| `recording/no-metrics-writes-nothing` | L2 RunTest | MetricsRoot has no new run JSONL |
| `recording/enabled-writes-run-start-end` | L2 RunTest (heavy) | JSONL has run_start then run_end |
| `recording/enabled-writes-leaf-events` | L2 RunTest (heavy) | leaf_start + leaf_end for executed leaf |
| `recording/phases/emits-tree-phases` | L2 RunTest | discover/generate/go_test phase events |
| `recording/phases/leaf-not-full-tree-wall` | L2 RunTest | leaf elapsed not full-tree clone |

## How to Run

```sh
doctest vet ./tests/metrics-record/
# default discovery: pure + flags + unlabeled recording (skips label: heavy)
doctest test ./tests/metrics-record/
# heavy recording leaves (nested prepare/go test via runner.RunTest)
doctest test --label heavy ./tests/metrics-record/...
# full suite
doctest test --label-all ./tests/metrics-record/...
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

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/metrics"
	"github.com/xhd2015/doctest/libdoc/runner"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

// Request selects one surface. Leaves set Op and related fields.
// All Ops are in-process (no product binary).
type Request struct {
	Op string // should_warn | format_warn | parse_flags | record_run

	// should_warn
	DefaultSuite bool
	Total        int
	Elapsed      time.Duration
	Threshold    time.Duration // 0 ⇒ metrics.DefaultSuiteWarnThreshold

	// should_warn multi-case (optional): when set, Run evaluates each case
	WarnCases []WarnCase

	// parse_flags
	Args []string

	// record_run
	MetricsRoot string // injectable root (t.TempDir); empty → Run creates one
	MetricsOn   bool
	Dir         string // fixture tree; empty → Run builds a 1-leaf pass tree
	LabelAll    bool
	LabelExprs  []string
}

type WarnCase struct {
	Name         string
	DefaultSuite bool
	Total        int
	Elapsed      time.Duration
	Threshold    time.Duration
	Want         bool
}

type Response struct {
	// pure warn
	ShouldWarn  bool
	Message     string
	WarnResults []bool // parallel to req.WarnCases

	// flags
	Opts       core.Options
	RemainArgs []string
	ParseErr   string

	// recording
	MetricsRoot string
	RunFiles    []string
	Decoded     []map[string]any // all events from first (or only) run file
	Stdout      string
	Stderr      string
	RunErr      string
	ExitCode    int
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}

	switch req.Op {
	case "should_warn":
		if len(req.WarnCases) > 0 {
			for _, c := range req.WarnCases {
				th := c.Threshold
				if th == 0 {
					th = metrics.DefaultSuiteWarnThreshold
				}
				got := metrics.ShouldWarnDefaultSuiteSlow(c.DefaultSuite, c.Total, c.Elapsed, th)
				resp.WarnResults = append(resp.WarnResults, got)
			}
			return resp, nil
		}
		th := req.Threshold
		if th == 0 {
			th = metrics.DefaultSuiteWarnThreshold
		}
		resp.ShouldWarn = metrics.ShouldWarnDefaultSuiteSlow(req.DefaultSuite, req.Total, req.Elapsed, th)
		return resp, nil

	case "format_warn":
		resp.Message = metrics.FormatDefaultSuiteSlowWarning()
		return resp, nil

	case "parse_flags":
		opts, remain, err := runner.ParseTestOptions(req.Args)
		if err != nil {
			resp.ParseErr = err.Error()
			return resp, nil
		}
		resp.Opts = opts
		resp.RemainArgs = remain
		return resp, nil

	case "record_run":
		root := req.MetricsRoot
		if root == "" {
			root = t.TempDir()
		}
		resp.MetricsRoot = root

		dir := req.Dir
		if dir == "" {
			dir = createOnePassTree(t)
		}

		before := listRunJSONL(t, root)

		opts := core.Options{
			Stderr:      &bytes.Buffer{},
			Stdout:      &bytes.Buffer{},
			RemoveTemp:  true,
			MetricsOn:   req.MetricsOn,
			MetricsRoot: root,
			LabelAll:    req.LabelAll,
			LabelExprs:  append([]string(nil), req.LabelExprs...),
		}
		err := runner.RunTest(dir, opts)
		if buf, ok := opts.Stdout.(*bytes.Buffer); ok {
			resp.Stdout = buf.String()
		}
		if buf, ok := opts.Stderr.(*bytes.Buffer); ok {
			resp.Stderr = buf.String()
		}
		if err != nil {
			resp.RunErr = err.Error()
		}

		after := listRunJSONL(t, root)
		resp.RunFiles = newPaths(before, after)
		if len(resp.RunFiles) > 0 {
			data, err := os.ReadFile(resp.RunFiles[0])
			if err != nil {
				return resp, err
			}
			resp.Decoded = decodeJSONL(data)
		}
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown Op %q", req.Op)
	}
}

func createOnePassTree(t *testing.T) string {
	t.Helper()
	// Process-shared fixture so recording leaves reuse gen/GOCACHE.
	return testtree.SharedPassFailTree(t, 1, 0)
}

func listRunJSONL(t *testing.T, metricsRoot string) map[string]struct{} {
	t.Helper()
	out := make(map[string]struct{})
	_ = filepath.Walk(metricsRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".jsonl") {
			out[path] = struct{}{}
		}
		return nil
	})
	return out
}

func newPaths(before, after map[string]struct{}) []string {
	var added []string
	for p := range after {
		if _, ok := before[p]; !ok {
			added = append(added, p)
		}
	}
	return added
}

func decodeJSONL(data []byte) []map[string]any {
	var decoded []map[string]any
	sc := bufio.NewScanner(bytes.NewReader(data))
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 4*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err == nil {
			decoded = append(decoded, m)
		} else {
			decoded = append(decoded, map[string]any{"_parse_error": err.Error(), "_raw": line})
		}
	}
	return decoded
}

func eventTypes(decoded []map[string]any) []string {
	var types []string
	for _, m := range decoded {
		if ty, _ := m["type"].(string); ty != "" {
			types = append(types, ty)
		}
	}
	return types
}

func hasEventType(decoded []map[string]any, typ string) bool {
	for _, m := range decoded {
		if m["type"] == typ {
			return true
		}
	}
	return false
}
```
