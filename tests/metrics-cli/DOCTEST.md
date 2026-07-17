# Metrics CLI — `doctest metrics` analyze (P3)

## Version
0.0.2

Classic-TDD specification for **Phase P3**: the `doctest metrics` user-facing
analyze CLI. Depends on P1 layout (`$MetricsRoot/doctest/metrics/<project_id>/runs/*.jsonl`)
and P2 recording shape (`run_start` / `leaf_*` / `run_end`, `mode.default_suite`).

These tests invoke the **project doctest binary** with fixture JSONL under an
injectable metrics root (`DOCTEST_METRICS_ROOT`). They do **not** run long live
suites to produce data. Expect **RED** until the implementer wires the
`metrics` subcommand and read/analyze path.

Out of scope: changing how tests are recorded (P2), review-perf skill content (P4).

# DSN (Domain Specific Notion)

### Participants

- **User / test harness** — runs `doctest metrics <subcmd>` from a project cwd.
- **doctest CLI** — binary under test; exposes top-level `metrics` with
  subcommands `path`, `last`, `top`, `summary`, `show`, `prune`, plus help.
- **Metrics root** — cache base directory. Production uses the user cache;
  tests inject via env `DOCTEST_METRICS_ROOT` (same override as P2 recording).
- **Project identity** — slug from `git remote get-url origin` in cwd (or
  `nogit_<hash>` fallback). Selects
  `$MetricsRoot/doctest/metrics/<project_id>/`.
- **Run store** — directory `runs/` of append-only JSONL files named
  `YYYY-MM-DD-HH-MM-SS-NN-<suffix>.jsonl` (UTC). Lexicographic filename order
  matches chronological order.
- **Event stream** — schema_version 1 objects: `run_start` (optional nested
  `mode.default_suite`, labels/mode flags), `leaf_start` (path, labels),
  `leaf_end` (path, elapsed_ns, result), `run_end` (wall_ns, passed, total, …).
- **Analyzer** — reads run files, ranks leaves by `elapsed_ns`, filters by
  suite mode / leaf labels, aggregates multi-run summaries, prunes old files.

### Behaviors

- **`metrics path`** — print the absolute project metrics directory
  (`…/doctest/metrics/<project_id>`) for cwd’s project identity under the
  resolved metrics root. Exit 0 even if the directory does not exist yet
  (path is the canonical location for future runs).
- **`metrics last`** — summarize the newest run file (by filename). Human
  stdout includes enough markers to identify the run (id/filename stem) and
  key counts from `run_end` when present. No runs → non-zero exit and a clear
  message on stderr (or combined output) mentioning no runs / not found.
- **`metrics top`** — rank `leaf_end` events by `elapsed_ns` descending from
  the selected run (`--run last` default, or `--run <id>` matching filename
  stem / basename without `.jsonl`). Flags:
  - `--n N` — limit rows (default implementation may use 10; tests set N).
  - `--unlabeled-only` — keep leaves whose labels are empty/absent
    (labels from matching `leaf_start` for the same path, else leaf_end).
  - `--default-only` — only consider runs with `mode.default_suite == true`
    (when selecting among runs / filtering the chosen run’s suite type).
  - `--json` — machine-readable array/object on stdout (valid JSON).
- **`metrics summary`** — aggregate the last N run files (`--last N`, default
  small positive N). Optional `--default-only` restricts to default-suite runs.
  `--json` emits machine-readable aggregate.
- **`metrics show [run-id]`** — dump one run: omit id ⇒ latest; with id ⇒ that
  file. Human output includes event type markers and/or leaf paths from the
  fixture. Unknown id ⇒ non-zero exit.
- **`metrics prune`** — delete oldest run files beyond default retention
  (**keep newest 30** files per project `runs/` directory, by filename order).
  When at or under retention, no files removed. Prints how many removed (or
  equivalent success text). Exit 0.
- **Help** — `doctest metrics --help` / `doctest help metrics` lists subcommands
  (`path`, `last`, `top`, `summary`, `show`, `prune`). Exit 0.
- **Unknown subcommand** — non-zero exit; stderr/stdout mentions unknown or usage.
- **Resolution** — project_id from cwd git origin (tests seed origin);
  metrics root from `DOCTEST_METRICS_ROOT` when set.

### Pipeline sketch

```
cwd + DOCTEST_METRICS_ROOT
  -> project_id (origin | nogit fallback)
  -> $root/doctest/metrics/<project_id>/runs/*.jsonl
  -> subcommand:
       path    -> print project metrics dir
       last    -> newest run summary
       top     -> rank leaf_end by elapsed_ns [filters] [--json]
       summary -> aggregate last N runs [--json]
       show    -> dump one run
       prune   -> keep newest 30, delete rest
```

## Decision Tree

```
tests/metrics-cli/
├── help/                                  [usage / errors]
│   ├── lists-subcommands/                 metrics --help lists path last top …
│   └── unknown-subcommand/                metrics nosuch → non-zero
├── path/                                  [canonical dir]
│   ├── prints-project-metrics-dir/        prints …/doctest/metrics/<id>
│   └── missing-dir-still-prints/          dir absent → still print path, exit 0
├── last/                                  [newest run]
│   ├── shows-latest-run/                  multi-file fixtures → newest summary
│   └── no-runs/                           empty store → non-zero + message
├── top/                                   [slow leaf ranking]
│   ├── ranks-slowest-leaves/              order by elapsed_ns desc
│   ├── n-limit/                           --n 2 truncates
│   ├── unlabeled-only/                    --unlabeled-only drops labeled leaves
│   ├── default-only/                      --default-only uses default-suite run
│   ├── run-select/                        --run <id> picks non-latest file
│   └── json-output/                       --json valid JSON with paths/times
├── summary/                               [multi-run aggregate]
│   ├── aggregates-last-n/                 --last 2 aggregates two newest
│   ├── default-only/                      skips non-default suite runs
│   └── json-output/                       --json valid JSON aggregate
├── show/                                  [single-run dump]
│   ├── latest-when-no-id/                 show → latest run content
│   ├── by-run-id/                         show <id> → that run
│   └── unknown-run-id/                    missing id → non-zero
└── prune/                                 [retention = 30]
    ├── removes-old-beyond-retention/      35 files → keep 30 newest
    └── under-retention-no-op/             3 files → still 3
```

## Test Index

| Leaf | Focus | Expected |
|------|--------|----------|
| `help/lists-subcommands` | help | exit 0; subcommand names present |
| `help/unknown-subcommand` | errors | exit ≠ 0 |
| `path/prints-project-metrics-dir` | path | stdout is/contains metrics project dir |
| `path/missing-dir-still-prints` | path | exit 0; path printed though dir missing |
| `last/shows-latest-run` | last | latest run id/stem + end stats markers |
| `last/no-runs` | empty | exit ≠ 0; no-run messaging |
| `top/ranks-slowest-leaves` | top | slowest path first |
| `top/n-limit` | top | only N rows |
| `top/unlabeled-only` | top | labeled leaf omitted |
| `top/default-only` | top | ranks from default-suite run |
| `top/run-select` | top | `--run` non-latest file |
| `top/json-output` | top | stdout JSON; includes slow path |
| `summary/aggregates-last-n` | summary | reflects N newest runs |
| `summary/default-only` | summary | excludes non-default runs |
| `summary/json-output` | summary | valid JSON |
| `show/latest-when-no-id` | show | latest fixture paths/events |
| `show/by-run-id` | show | selected older run |
| `show/unknown-run-id` | show | exit ≠ 0 |
| `prune/removes-old-beyond-retention` | prune | 35 → 30 files |
| `prune/under-retention-no-op` | prune | count unchanged |

## How to Run

```sh
doctest vet ./tests/metrics-cli/
doctest test ./tests/metrics-cli/                 # RED until metrics CLI wired
doctest test ./tests/metrics-cli/help/...
doctest test ./tests/metrics-cli/path/...
doctest test ./tests/metrics-cli/last/...
doctest test ./tests/metrics-cli/top/...
doctest test ./tests/metrics-cli/summary/...
doctest test ./tests/metrics-cli/show/...
doctest test ./tests/metrics-cli/prune/...
```

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/metrics"
)

// FixtureOrigin is the git remote used so project_id is stable across leaves.
const FixtureOrigin = "https://github.com/xhd2015/doctest-metrics-fixture.git"

// DefaultRunRetention is the P3 prune keep-count (newest N run files).
const DefaultRunRetention = 30

// Request drives one `doctest metrics …` invocation against fixtures.
type Request struct {
	Args        []string // argv after binary, e.g. ["metrics","top","--n","3"]
	Env         []string // extra env (DOCTEST_METRICS_ROOT set by Run if MetricsRoot set)
	WorkDir     string   // subprocess cwd (git project)
	MetricsRoot string   // injectable cache root
	ProjectID   string   // expected slug (usually from origin)
	Timeout     time.Duration
	Bin         string

	// SnapshotRunFilesAfter: if true, Response.RunFiles lists *.jsonl basenames after command.
	SnapshotRunFilesAfter bool
}

type Response struct {
	ExitCode    int
	Stdout      string
	Stderr      string
	Err         error
	MetricsRoot string
	ProjectID   string
	RunsDir     string
	RunFiles    []string // basenames, sorted
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req.Bin == "" {
		return nil, fmt.Errorf("req.Bin is not set")
	}
	timeout := req.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	env := append([]string{}, os.Environ()...)
	env = append(env, req.Env...)
	if req.MetricsRoot != "" {
		env = append(env, "DOCTEST_METRICS_ROOT="+req.MetricsRoot)
	}

	cmd := exec.CommandContext(ctx, req.Bin, req.Args...)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	}
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()

	resp := &Response{
		Stdout:      stdout.String(),
		Stderr:      stderr.String(),
		Err:         runErr,
		MetricsRoot: req.MetricsRoot,
		ProjectID:   req.ProjectID,
		ExitCode:    0,
	}
	if req.ProjectID != "" && req.MetricsRoot != "" {
		resp.RunsDir = filepath.Join(req.MetricsRoot, "doctest", "metrics", req.ProjectID, "runs")
	}
	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			resp.ExitCode = exitErr.ExitCode()
			runErr = nil
		} else if ctx.Err() != nil {
			return resp, ctx.Err()
		} else {
			return resp, runErr
		}
	}
	if req.SnapshotRunFilesAfter && resp.RunsDir != "" {
		resp.RunFiles = listJSONLBasenames(resp.RunsDir)
	}
	return resp, nil
}

func listJSONLBasenames(dir string) []string {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".jsonl") {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

// compile-time touch so metrics import stays available for helpers in SETUP.
var _ = metrics.SchemaVersion
```
