# Metrics CLI — analyze recorded metrics (in-process + sparse e2e)

## Version
0.0.2

**Layer model (coverage backfill):**

| Layer | Share | Where |
|-------|-------|--------|
| **L2 doctest in-process** | **mass (≥80%)** | `path/`, `last/`, `top/`, `summary/`, `show/`, `prune/` — harness `Run` calls `libdoc/metrics` analyze APIs with fixture JSONL under `t.TempDir()` |
| **L3 doctest e2e** | **sparse (≤3)** | `e2e/` — real `doctest metrics` binary for help/argv wiring only; `label: heavy` |

Default discovery runs L2 only (unlabeled). Use `--label heavy` for binary smokes.

Depends on P1 layout (`$MetricsRoot/doctest/metrics/<project_id>/runs/*.jsonl`) and
P2 event shapes (`run_start` / `leaf_*` / `run_end`, `mode.default_suite`).
Fixtures are written under an injectable metrics root — **no long live suites**.

Out of scope: changing how tests are recorded (`tests/metrics-record`),
review-perf skill content, metrics-foundation APIs.

# DSN (Domain Specific Notion)

### Participants

- **Harness** — default path: calls `libdoc/metrics` analyze functions in-process
  (`ListRunFiles`, `RankLeaves`, `SelectRun`, `AggregateRuns`, `PruneRuns`,
  `SummarizeRun`, `ProjectMetricsDir`, …) against fixture JSONL.
- **e2e binary** — sparse L3 leaves spawn `doctest metrics` for help / unknown
  subcommand wiring only.
- **Metrics root** — cache base directory; tests inject via `Request.MetricsRoot`
  (CLI env `DOCTEST_METRICS_ROOT` only in e2e).
- **Project identity** — slug from fixture origin (or `ProjectIDForDir`); selects
  `$MetricsRoot/doctest/metrics/<project_id>/`.
- **Run store** — directory `runs/` of append-only JSONL files named
  `YYYY-MM-DD-HH-MM-SS-NN-<suffix>.jsonl` (UTC). Lexicographic filename order
  matches chronological order.
- **Event stream** — schema_version 1 objects: `run_start`, `leaf_start`,
  `leaf_end`, `run_end`.
- **Analyzer** — ranks leaves by `elapsed_ns`, filters by suite mode / labels,
  aggregates multi-run summaries, prunes old files.

### Behaviors

- **`path`** — resolve absolute project metrics directory
  (`…/doctest/metrics/<project_id>`). Exit 0 even if the directory does not
  exist yet.
- **`last`** — summarize the newest run file. No runs → non-zero + clear message.
- **`top`** — rank `leaf_end` by `elapsed_ns` desc. Flags: `--n`,
  `--unlabeled-only`, `--default-only`, `--run`, `--json`.
- **`summary`** — aggregate last N runs (`--last N`, `--default-only`, `--json`).
- **`show [run-id]`** — dump one run (latest if id omitted). Unknown id → non-zero.
- **`prune`** — keep newest **30** run files; delete older. Exit 0.
- **Help / unknown** — L3 e2e only: lists subcommands; unknown → non-zero.

### Pipeline sketch

```
# L2 (default)
MetricsRoot + ProjectID + fixture JSONL
  -> ProjectRunsDir / ListRunFiles
  -> subcommand via metrics.* APIs
       path    -> ProjectMetricsDir
       last    -> NewestRun + SummarizeRun
       top     -> SelectRun + ExtractLeaves + RankLeaves
       summary -> LastNRuns + AggregateRuns
       show    -> FindRunByID / NewestRun + ReadEvents
       prune   -> PruneRuns(keep=30)

# L3 (e2e/, label: heavy)
req.Bin metrics <subcmd> + DOCTEST_METRICS_ROOT
```

## Decision Tree

```
tests/metrics-cli/
├── path/                                  [L2 ProjectMetricsDir]
│   ├── prints-project-metrics-dir/
│   └── missing-dir-still-prints/
├── last/                                  [L2 NewestRun + SummarizeRun]
│   ├── shows-latest-run/
│   └── no-runs/
├── top/                                   [L2 RankLeaves / SelectRun]
│   ├── ranks-slowest-leaves/
│   ├── n-limit/
│   ├── unlabeled-only/
│   ├── default-only/
│   ├── run-select/
│   └── json-output/
├── summary/                               [L2 LastNRuns + AggregateRuns]
│   ├── aggregates-last-n/
│   ├── default-only/
│   └── json-output/
├── show/                                  [L2 ReadEvents / FindRunByID]
│   ├── latest-when-no-id/
│   ├── by-run-id/
│   └── unknown-run-id/
├── prune/                                 [L2 PruneRuns keep=30]
│   ├── removes-old-beyond-retention/
│   └── under-retention-no-op/
└── e2e/                                   [L3 binary, label: heavy]
    ├── help-lists-subcommands/
    └── unknown-subcommand/
```

## Test Index

| Leaf | Layer | Expected |
|------|--------|----------|
| `path/prints-project-metrics-dir` | L2 | metrics project dir path |
| `path/missing-dir-still-prints` | L2 | exit 0; path printed though dir missing |
| `last/shows-latest-run` | L2 | newest run id/stem + stats |
| `last/no-runs` | L2 | exit ≠ 0; no-run messaging |
| `top/ranks-slowest-leaves` | L2 | slowest path first |
| `top/n-limit` | L2 | only N rows |
| `top/unlabeled-only` | L2 | labeled leaf omitted |
| `top/default-only` | L2 | ranks from default-suite run |
| `top/run-select` | L2 | non-latest file |
| `top/json-output` | L2 | JSON with paths/times |
| `summary/aggregates-last-n` | L2 | N newest runs |
| `summary/default-only` | L2 | excludes non-default runs |
| `summary/json-output` | L2 | valid JSON |
| `show/latest-when-no-id` | L2 | latest fixture paths/events |
| `show/by-run-id` | L2 | selected older run |
| `show/unknown-run-id` | L2 | exit ≠ 0 |
| `prune/removes-old-beyond-retention` | L2 | 35 → 30 files |
| `prune/under-retention-no-op` | L2 | count unchanged |
| `e2e/help-lists-subcommands` | L3 heavy | exit 0; subcommand names |
| `e2e/unknown-subcommand` | L3 heavy | exit ≠ 0 |

## How to Run

```sh
doctest vet ./tests/metrics-cli/
# default discovery: L2 in-process mass (skips label: heavy)
doctest test ./tests/metrics-cli/
# sparse binary smokes
doctest test --label heavy ./tests/metrics-cli/...
# full suite
doctest test --label-all ./tests/metrics-cli/...
```

```go
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/metrics"
)

// FixtureOrigin is the git remote used so project_id is stable across leaves.
const FixtureOrigin = "https://github.com/xhd2015/doctest-metrics-fixture.git"

// DefaultRunRetention is the prune keep-count (newest N run files).
const DefaultRunRetention = 30

// Request drives one metrics analyze scenario against fixtures.
// Default Mode is in-process (libdoc/metrics). Set UseCLI for L3 e2e only.
type Request struct {
	Args        []string // e.g. ["metrics","top","--n","3"] — parsed for L2 dispatch
	Env         []string // e2e only
	WorkDir     string   // project cwd (git); used for ProjectIDForDir when needed
	MetricsRoot string   // injectable cache root
	ProjectID   string   // expected slug (usually from origin)
	Timeout     time.Duration
	Bin         string // e2e only: product binary

	// UseCLI: when true, spawn req.Bin (L3 e2e). Default false = in-process APIs.
	UseCLI bool

	// SnapshotRunFilesAfter: if true, Response.RunFiles lists *.jsonl basenames after op.
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
	if req.UseCLI {
		return runE2E(t, req)
	}
	return runInProcess(t, req)
}

func runE2E(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req.Bin == "" {
		return nil, fmt.Errorf("UseCLI requires req.Bin")
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

func runInProcess(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{
		MetricsRoot: req.MetricsRoot,
		ProjectID:   req.ProjectID,
		ExitCode:    0,
	}
	if req.ProjectID != "" && req.MetricsRoot != "" {
		resp.RunsDir = metrics.ProjectRunsDir(req.MetricsRoot, req.ProjectID)
	}

	args := append([]string(nil), req.Args...)
	if len(args) > 0 && args[0] == "metrics" {
		args = args[1:]
	}
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		// Help is covered by L3 e2e; still provide a usable message if invoked.
		resp.Stdout = "Usage: doctest metrics <subcommand>\npath last top summary show prune\n"
		return resp, nil
	}

	sub := args[0]
	rest := args[1:]
	var opErr error
	switch sub {
	case "path":
		opErr = inProcessPath(req, resp)
	case "last":
		opErr = inProcessLast(req, resp)
	case "top":
		opErr = inProcessTop(req, resp, rest)
	case "summary":
		opErr = inProcessSummary(req, resp, rest)
	case "show":
		opErr = inProcessShow(req, resp, rest)
	case "prune":
		opErr = inProcessPrune(req, resp)
	default:
		opErr = fmt.Errorf("unknown metrics subcommand: %s", sub)
	}
	if opErr != nil {
		resp.ExitCode = 1
		resp.Stderr = opErr.Error() + "\n"
		resp.Err = opErr
	}
	if req.SnapshotRunFilesAfter {
		runs := resp.RunsDir
		if runs == "" && req.MetricsRoot != "" && req.ProjectID != "" {
			runs = metrics.ProjectRunsDir(req.MetricsRoot, req.ProjectID)
			resp.RunsDir = runs
		}
		if runs != "" {
			resp.RunFiles = listJSONLBasenames(runs)
		}
	}
	return resp, nil
}

func resolveProject(req *Request) (root, projectID, metricsDir, runsDir string, err error) {
	root = metrics.ResolveMetricsRoot(req.MetricsRoot)
	if root == "" {
		return "", "", "", "", fmt.Errorf("cannot resolve metrics root")
	}
	projectID = req.ProjectID
	if projectID == "" && req.WorkDir != "" {
		projectID = metrics.ProjectIDForDir(req.WorkDir)
	}
	if projectID == "" {
		return "", "", "", "", fmt.Errorf("cannot resolve project id")
	}
	metricsDir = metrics.ProjectMetricsDir(root, projectID)
	runsDir = metrics.ProjectRunsDir(root, projectID)
	return root, projectID, metricsDir, runsDir, nil
}

func inProcessPath(req *Request, resp *Response) error {
	_, pid, mdir, _, err := resolveProject(req)
	if err != nil {
		return err
	}
	resp.ProjectID = pid
	abs, err := filepath.Abs(mdir)
	if err != nil {
		abs = mdir
	}
	resp.Stdout = abs + "\n"
	return nil
}

func inProcessLast(req *Request, resp *Response) error {
	_, _, _, runsDir, err := resolveProject(req)
	if err != nil {
		return err
	}
	resp.RunsDir = runsDir
	files, err := metrics.ListRunFiles(runsDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no runs found in %s", runsDir)
	}
	rf := metrics.NewestRun(files)
	evs, err := metrics.ReadEvents(rf.Path)
	if err != nil {
		return err
	}
	s := metrics.SummarizeRun(*rf, evs)
	resp.Stdout = formatRunSummary(s)
	return nil
}

func formatRunSummary(s metrics.RunSummary) string {
	var b strings.Builder
	fmt.Fprintf(&b, "run_id: %s\n", s.RunID)
	fmt.Fprintf(&b, "file: %s\n", s.File)
	fmt.Fprintf(&b, "default_suite: %v\n", s.DefaultSuite)
	if s.HasRunEnd {
		fmt.Fprintf(&b, "passed: %d  total: %d  skipped: %d\n", s.Passed, s.Total, s.Skipped)
		fmt.Fprintf(&b, "wall: %s\n", metrics.FormatDurationNS(s.WallNs))
		fmt.Fprintf(&b, "exit_ok: %v\n", s.ExitOK)
	}
	fmt.Fprintf(&b, "leaf_count: %d\n", s.LeafCount)
	if len(s.Slowest) > 0 {
		b.WriteString("slowest:\n")
		for _, row := range s.Slowest {
			fmt.Fprintf(&b, "  %s  %s  %s\n", row.Path, metrics.FormatDurationNS(row.ElapsedNs), row.Result)
		}
	}
	return b.String()
}

func parseTopFlags(args []string) (n int, unlabeledOnly, defaultOnly, asJSON bool, runSel string, err error) {
	n = metrics.DefaultTopN
	runSel = "last"
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--unlabeled-only":
			unlabeledOnly = true
		case a == "--default-only":
			defaultOnly = true
		case a == "--json":
			asJSON = true
		case a == "--n" || a == "-n":
			if i+1 >= len(args) {
				return 0, false, false, false, "", fmt.Errorf("--n requires a value")
			}
			i++
			v, e := strconv.Atoi(args[i])
			if e != nil || v < 0 {
				return 0, false, false, false, "", fmt.Errorf("invalid --n: %s", args[i])
			}
			n = v
		case strings.HasPrefix(a, "--n="):
			v, e := strconv.Atoi(strings.TrimPrefix(a, "--n="))
			if e != nil || v < 0 {
				return 0, false, false, false, "", fmt.Errorf("invalid --n: %s", a)
			}
			n = v
		case a == "--run":
			if i+1 >= len(args) {
				return 0, false, false, false, "", fmt.Errorf("--run requires a value")
			}
			i++
			runSel = args[i]
		case strings.HasPrefix(a, "--run="):
			runSel = strings.TrimPrefix(a, "--run=")
		case a == "-h" || a == "--help":
			// ignore
		default:
			return 0, false, false, false, "", fmt.Errorf("unknown flag for metrics top: %s", a)
		}
	}
	return n, unlabeledOnly, defaultOnly, asJSON, runSel, nil
}

func inProcessTop(req *Request, resp *Response, args []string) error {
	n, unlabeledOnly, defaultOnly, asJSON, runSel, err := parseTopFlags(args)
	if err != nil {
		return err
	}
	_, _, _, runsDir, err := resolveProject(req)
	if err != nil {
		return err
	}
	resp.RunsDir = runsDir
	files, err := metrics.ListRunFiles(runsDir)
	if err != nil {
		return err
	}
	rf, err := metrics.SelectRun(files, runSel, defaultOnly)
	if err != nil {
		return err
	}
	evs, err := metrics.ReadEvents(rf.Path)
	if err != nil {
		return err
	}
	rows := metrics.RankLeaves(metrics.ExtractLeaves(evs), unlabeledOnly, n)
	if asJSON {
		data, err := json.MarshalIndent(rows, "", "  ")
		if err != nil {
			return err
		}
		resp.Stdout = string(data) + "\n"
		return nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "run_id: %s\n", rf.ID)
	for i, row := range rows {
		fmt.Fprintf(&b, "%d. %s  %s", i+1, row.Path, metrics.FormatDurationNS(row.ElapsedNs))
		if row.Result != "" {
			fmt.Fprintf(&b, "  %s", row.Result)
		}
		if len(row.Labels) > 0 {
			fmt.Fprintf(&b, "  labels=%s", strings.Join(row.Labels, ","))
		}
		b.WriteByte('\n')
	}
	resp.Stdout = b.String()
	return nil
}

func parseSummaryFlags(args []string) (lastN int, defaultOnly, asJSON bool, err error) {
	lastN = metrics.DefaultSummaryLast
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--default-only":
			defaultOnly = true
		case a == "--json":
			asJSON = true
		case a == "--last":
			if i+1 >= len(args) {
				return 0, false, false, fmt.Errorf("--last requires a value")
			}
			i++
			v, e := strconv.Atoi(args[i])
			if e != nil || v < 0 {
				return 0, false, false, fmt.Errorf("invalid --last: %s", args[i])
			}
			lastN = v
		case strings.HasPrefix(a, "--last="):
			v, e := strconv.Atoi(strings.TrimPrefix(a, "--last="))
			if e != nil || v < 0 {
				return 0, false, false, fmt.Errorf("invalid --last: %s", a)
			}
			lastN = v
		case a == "-h" || a == "--help":
			// ignore
		default:
			return 0, false, false, fmt.Errorf("unknown flag for metrics summary: %s", a)
		}
	}
	return lastN, defaultOnly, asJSON, nil
}

func inProcessSummary(req *Request, resp *Response, args []string) error {
	lastN, defaultOnly, asJSON, err := parseSummaryFlags(args)
	if err != nil {
		return err
	}
	_, _, _, runsDir, err := resolveProject(req)
	if err != nil {
		return err
	}
	resp.RunsDir = runsDir
	files, err := metrics.ListRunFiles(runsDir)
	if err != nil {
		return err
	}
	selected, err := metrics.LastNRuns(files, lastN, defaultOnly)
	if err != nil {
		return err
	}
	if len(selected) == 0 {
		return fmt.Errorf("no runs found")
	}
	agg, err := metrics.AggregateRuns(selected)
	if err != nil {
		return err
	}
	agg.DefaultOnly = defaultOnly
	if asJSON {
		data, err := json.MarshalIndent(agg, "", "  ")
		if err != nil {
			return err
		}
		resp.Stdout = string(data) + "\n"
		return nil
	}
	var b strings.Builder
	fmt.Fprintf(&b, "runs: %d\n", agg.Runs)
	fmt.Fprintf(&b, "run_ids: %s\n", strings.Join(agg.RunIDs, ", "))
	fmt.Fprintf(&b, "passed: %d  total: %d  skipped: %d\n", agg.Passed, agg.Total, agg.Skipped)
	fmt.Fprintf(&b, "wall: %s\n", metrics.FormatDurationNS(agg.WallNs))
	fmt.Fprintf(&b, "leaf_count: %d\n", agg.LeafCount)
	if defaultOnly {
		b.WriteString("default_only: true\n")
	}
	resp.Stdout = b.String()
	return nil
}

func inProcessShow(req *Request, resp *Response, args []string) error {
	runID := ""
	for _, a := range args {
		if a == "-h" || a == "--help" {
			continue
		}
		if strings.HasPrefix(a, "-") {
			return fmt.Errorf("unknown flag for metrics show: %s", a)
		}
		if runID == "" {
			runID = a
		} else {
			return fmt.Errorf("metrics show accepts at most one run-id")
		}
	}
	_, _, _, runsDir, err := resolveProject(req)
	if err != nil {
		return err
	}
	resp.RunsDir = runsDir
	files, err := metrics.ListRunFiles(runsDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no runs found")
	}
	var rf *metrics.RunFile
	if runID == "" {
		rf = metrics.NewestRun(files)
	} else {
		rf = metrics.FindRunByID(files, runID)
		if rf == nil {
			return fmt.Errorf("run not found: %s", runID)
		}
	}
	evs, err := metrics.ReadEvents(rf.Path)
	if err != nil {
		return err
	}
	var b strings.Builder
	fmt.Fprintf(&b, "run_id: %s\n", rf.ID)
	fmt.Fprintf(&b, "file: %s\n", rf.Name)
	for _, ev := range evs {
		typ, _ := ev["type"].(string)
		switch typ {
		case "run_start":
			fmt.Fprintf(&b, "event: run_start run_id=%v\n", ev["run_id"])
			if mode, ok := ev["mode"].(map[string]any); ok {
				fmt.Fprintf(&b, "  mode.default_suite=%v\n", mode["default_suite"])
			}
		case "leaf_start":
			fmt.Fprintf(&b, "event: leaf_start path=%v labels=%v\n", ev["path"], ev["labels"])
		case "leaf_end":
			fmt.Fprintf(&b, "event: leaf_end path=%v elapsed_ns=%v result=%v\n", ev["path"], ev["elapsed_ns"], ev["result"])
		case "run_end":
			fmt.Fprintf(&b, "event: run_end passed=%v total=%v wall_ns=%v exit_ok=%v\n",
				ev["passed"], ev["total"], ev["wall_ns"], ev["exit_ok"])
		default:
			raw, _ := json.Marshal(ev)
			fmt.Fprintf(&b, "event: %s\n", string(raw))
		}
	}
	resp.Stdout = b.String()
	return nil
}

func inProcessPrune(req *Request, resp *Response) error {
	_, _, _, runsDir, err := resolveProject(req)
	if err != nil {
		return err
	}
	resp.RunsDir = runsDir
	removed, err := metrics.PruneRuns(runsDir, metrics.DefaultRunRetention)
	if err != nil {
		return err
	}
	resp.Stdout = fmt.Sprintf("prune: removed %d file(s); retention=%d\n", removed, metrics.DefaultRunRetention)
	return nil
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
