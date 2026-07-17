# Metrics Record — doctest test JSONL + default-suite WARNING (P2)

## Version
0.0.2

Classic-TDD specification for **Phase P2**: wire recording from `doctest test`
and the non-fatal default-suite slow WARNING. Depends on P1 package
`github.com/xhd2015/doctest/libdoc/metrics` (writer, paths, project id).

These tests exercise pure warn helpers and suite-level recording options. They
do **not** implement production wiring — expect **RED** until the implementer
adds helpers, `core.Options` fields, runner/CLI flags, and suite hooks.

Out of scope: `doctest metrics` analyze (P3), review-perf skill docs (P4),
env silence for WARNING, fail-on-slow-suite.

# DSN (Domain Specific Notion)

### Participants

- **Caller** — `doctest test` (CLI or package entry) that owns one suite wall clock.
- **Warn helpers** — pure functions deciding whether to emit the default-suite
  slow WARNING and formatting its fixed message (no real 3-minute sleep in tests).
- **Metrics options** — `NoMetrics` (CLI `--no-metrics`) and injectable
  `MetricsRoot` (test cache root; production uses the user cache layout from P1).
- **Run recorder** — opens one exclusive JSONL run file under
  `$MetricsRoot/doctest/metrics/<project_id>/runs/`, writes `run_start` /
  leaf events / `run_end` via the P1 buffered writer, closes on suite end.
- **Default suite** — discovery with `!LabelAll && len(LabelExprs)==0`.
- **Stderr** — receives non-fatal `WARNING:` after the result summary when the
  warn predicate fires; process exit is never driven solely by this warning.

### Behaviors

- **Should warn** — true only when default suite, `total > 0`, and
  `elapsed > DefaultSuiteWarnThreshold` (3 minutes). False when labeled /
  `--label-all`, when elapsed ≤ threshold, or when total is 0.
- **Format warning** — fixed prose mentioning `WARNING:`, default suite speed
  (3 minutes), `skill:doctest-review-perf`, and `doctest skill review-perf --show`.
- **Opt-out** — `--no-metrics` / `NoMetrics` skips creating any run file under
  MetricsRoot (and does not open a writer).
- **Record** — with metrics on: write `run_start` (project_id, cwd, argv/flags,
  git branch/commit without dirty, session_id, `mode.default_suite`,
  `schema_version` 1), leaf_start/leaf_end for executed leaves when practical,
  then `run_end` (wall_ns, passed, total, skipped, exit_ok, warnings including
  `default_suite_slow` when the warn predicate fired).
- **Emit WARNING** — after suite finish (after result summary preferred); yellow
  when color is enabled (color itself is not asserted in this tree).

### Pipeline sketch

```
doctest test [flags] <dir>
  -> parse --no-metrics / labels -> Options
  -> if !NoMetrics: open run JSONL under MetricsRoot (or user cache)
  -> run_start
  -> execute leaves -> leaf_start / leaf_end*
  -> run_end (+ warnings)
  -> PrintResultSummary
  -> if ShouldWarnDefaultSuiteSlow: FormatDefaultSuiteSlowWarning -> stderr
```

## Decision Tree

```
tests/metrics-record/
├── warn-predicate/                      [pure ShouldWarnDefaultSuiteSlow]
│   ├── fires-default-over-threshold/    default + total>0 + elapsed>3m → true
│   ├── no-fire-elapsed-at-or-under/     default + elapsed≤3m → false
│   ├── no-fire-label-all/               LabelAll + slow → false
│   ├── no-fire-label-exprs/             LabelExprs set + slow → false
│   └── no-fire-total-zero/              total==0 + slow → false
├── warning-message/                     [FormatDefaultSuiteSlowWarning]
│   └── required-phrases/                WARNING:, skill, review-perf --show, 3 minutes
├── flags/                               [CLI / ParseTestOptions]
│   ├── default-metrics-on/              omit --no-metrics → NoMetrics=false
│   └── no-metrics-sets-opt-out/         --no-metrics → NoMetrics=true
└── recording/                           [suite run under injectable MetricsRoot]
    ├── no-metrics-writes-nothing/       --no-metrics → no new *.jsonl under root
    ├── enabled-writes-run-start-end/    metrics on → run_start + run_end present
    └── enabled-writes-leaf-events/      metrics on + 1 leaf → leaf_start/end present
```

## Test Index

| Leaf | Focus | Expected |
|------|--------|----------|
| `warn-predicate/fires-default-over-threshold` | pure | `ShouldWarn...` true |
| `warn-predicate/no-fire-elapsed-at-or-under` | pure | false at 3m and under 3m |
| `warn-predicate/no-fire-label-all` | pure | false when not default (label-all) |
| `warn-predicate/no-fire-label-exprs` | pure | false when LabelExprs non-empty |
| `warn-predicate/no-fire-total-zero` | pure | false when total=0 |
| `warning-message/required-phrases` | pure | fixed message substrings |
| `flags/default-metrics-on` | parse | `NoMetrics == false` by default |
| `flags/no-metrics-sets-opt-out` | parse | `--no-metrics` → `NoMetrics == true` |
| `recording/no-metrics-writes-nothing` | integration | MetricsRoot has no new run JSONL |
| `recording/enabled-writes-run-start-end` | integration | JSONL has run_start then run_end |
| `recording/enabled-writes-leaf-events` | integration | leaf_start + leaf_end for executed leaf |

## How to Run

```sh
doctest vet ./tests/metrics-record/
doctest test ./tests/metrics-record/                 # RED until P2 wired
doctest test ./tests/metrics-record/warn-predicate/...
doctest test ./tests/metrics-record/warning-message/...
doctest test ./tests/metrics-record/flags/...
doctest test ./tests/metrics-record/recording/...
```

```go
import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/metrics"
	"github.com/xhd2015/doctest/libdoc/runner"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

// Request selects one P2 surface. Leaves set Op and related fields.
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
	NoMetrics   bool
	Dir         string // fixture tree; empty → Run builds a 1-leaf pass tree
	LabelAll    bool
	LabelExprs  []string
	ExtraArgs   []string // appended to `test` argv when using CLI path
	UseCLI      bool     // if true, invoke doctest binary; else package RunTest
	Bin         string   // set by Setup when UseCLI
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
	ShouldWarn bool
	Message    string
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

func Run(t *testing.T, req *Request) (*Response, error) {
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
		// Implementer exports ParseTestOptions (or TestExported_ParseTestOptions).
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

		if req.UseCLI {
			if req.Bin == "" {
				return nil, fmt.Errorf("record_run UseCLI requires req.Bin")
			}
			args := []string{"test"}
			if req.NoMetrics {
				args = append(args, "--no-metrics")
			}
			if req.LabelAll {
				args = append(args, "--label-all")
			}
			for _, e := range req.LabelExprs {
				args = append(args, "--label", e)
			}
			args = append(args, req.ExtraArgs...)
			args = append(args, dir)

			cmd := exec.Command(req.Bin, args...)
			cmd.Env = append(os.Environ(), "DOCTEST_METRICS_ROOT="+root)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()
			resp.Stdout = stdout.String()
			resp.Stderr = stderr.String()
			if err != nil {
				resp.RunErr = err.Error()
				if ee, ok := err.(*exec.ExitError); ok {
					resp.ExitCode = ee.ExitCode()
				} else {
					resp.ExitCode = -1
				}
			}
		} else {
			opts := core.Options{
				Stderr:      &bytes.Buffer{},
				Stdout:      &bytes.Buffer{},
				RemoveTemp:  true,
				NoMetrics:   req.NoMetrics,
				MetricsRoot: root,
				LabelAll:    req.LabelAll,
				LabelExprs:  append([]string(nil), req.LabelExprs...),
			}
			// Implementer provides suite entry that records metrics + optional WARNING.
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
	tmp := t.TempDir()
	testtree.WritePassFailTree(t, tmp, 1, 0)
	return tmp
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
