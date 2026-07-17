# Metrics Integration Polish — end-to-end wiring (P5)

## Version
0.0.2

Coverage-backfill / mixed-mode specification for **Phase P5**: one coherent
product story across recording, metrics CLI, help surfaces, and the
review-perf WARNING ↔ skill cross-link.

Depends on P1–P4 (`libdoc/metrics`, suite recording, `doctest metrics`,
`doctest skill review-perf`). Prefer **GREEN** when wiring already works;
**RED** is OK for help gaps (e.g. `test --help` missing `--no-metrics`).

Out of scope: phase timers, fsync, dirty git, fail-on-slow, new analyze flags.

# DSN (Domain Specific Notion)

### Participants

- **User / harness** — runs a tiny `doctest test` suite, then analyzes metrics
  or reads help/skill text.
- **doctest CLI** — binary under test (`test` records; `metrics` analyzes;
  `--help` / `test --help` documents surfaces; `skill review-perf --show`
  documents the perf skill).
- **Tiny suite** — one-leaf pass fixture written via `testtree` under a temp
  dir (no multi-minute wall clock).
- **Metrics root** — injectable cache base via `MetricsRoot` / env
  `DOCTEST_METRICS_ROOT`. Layout:
  `$root/doctest/metrics/<project_id>/runs/*.jsonl`.
- **Project identity** — derived from the suite directory (git origin or
  `nogit_*` fallback). Smoke leaves keep **record cwd and metrics cwd** on
  the same fixture path so project_id matches.
- **WARNING banner** — fixed string from
  `metrics.FormatDefaultSuiteSlowWarning()` (phrases: `WARNING:`,
  `skill:doctest-review-perf`, `review-perf --show`, `3 minutes`).
- **review-perf skill** — embedded doc shown by `skill review-perf --show`;
  must mention the same guidance phrases the WARNING recommends.

### Behaviors

- **Smoke record → analyze** — with metrics on and injectable MetricsRoot:
  running a tiny suite creates ≥1 run JSONL; then `doctest metrics last`
  and/or `doctest metrics top` exit 0 and surface that run (run stem / leaf
  path / totals) under the same MetricsRoot + project cwd.
- **Top-level help** — `doctest --help` (or bare usage) lists the `metrics`
  command among product surfaces.
- **Test help** — `doctest test --help` documents `--no-metrics` (opt-out of
  suite metrics recording).
- **Skill ↔ WARNING alignment** — every required phrase from
  `FormatDefaultSuiteSlowWarning` appears in the review-perf skill body (so
  the banner does not point at absent guidance). P4 already covers skill show
  and pure warn formatting separately; this leaf is a thin cross-link.

### Pipeline sketch

```
# smoke
create 1-leaf pass tree (temp)
  -> RunTest | doctest test  with MetricsRoot / DOCTEST_METRICS_ROOT
  -> runs/*.jsonl appears
  -> doctest metrics last|top (same root + cwd) -> exit 0 + run evidence

# help
doctest --help -> mentions metrics
doctest test --help -> mentions --no-metrics

# alignment
FormatDefaultSuiteSlowWarning() phrases ⊆ skill review-perf --show
```

## Decision Tree

```
tests/metrics-integration/
├── smoke/                                      [record → analyze e2e]
│   ├── record-then-last/                       JSONL then metrics last
│   └── record-then-top/                        JSONL then metrics top
├── help/                                       [CLI discovery surfaces]
│   ├── top-level-lists-metrics/                doctest --help includes metrics
│   └── test-usage-mentions-no-metrics/         test --help includes --no-metrics
└── skill-warning-alignment/                    [banner ↔ skill phrases]
    └── warn-phrases-in-skill-show/             FormatDefaultSuiteSlowWarning ⊆ skill
```

## Test Index

| Leaf | Focus | Expected |
|------|--------|----------|
| `smoke/record-then-last` | e2e | suite writes JSONL; `metrics last` exit 0 + run evidence |
| `smoke/record-then-top` | e2e | suite writes JSONL; `metrics top` exit 0 + leaf path |
| `help/top-level-lists-metrics` | help | exit 0; stdout contains `metrics` |
| `help/test-usage-mentions-no-metrics` | help | exit 0; stdout contains `--no-metrics` |
| `skill-warning-alignment/warn-phrases-in-skill-show` | cross-link | skill show contains WARNING required phrases |

## How to Run

```sh
doctest vet ./tests/metrics-integration/
doctest test ./tests/metrics-integration/
doctest test ./tests/metrics-integration/smoke/...
doctest test ./tests/metrics-integration/help/...
doctest test ./tests/metrics-integration/skill-warning-alignment/...
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
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/metrics"
	"github.com/xhd2015/doctest/libdoc/runner"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

// FixtureOrigin seeds the smoke fixture git remote so project_id is stable and
// identical for RunTest recording and the metrics CLI (both resolve from dir/cwd).
const FixtureOrigin = "https://github.com/xhd2015/doctest-metrics-integration.git"

// Request selects one P5 surface. Leaves set Op and related fields.
type Request struct {
	Op string // smoke | help | align_skill_warn

	// Shared CLI
	Bin     string
	Timeout time.Duration
	Args    []string // help / skill argv after binary

	// smoke
	MetricsRoot string // injectable; empty → Run creates temp
	FixtureDir  string // tiny suite root; empty → Run builds 1-leaf pass tree
	// AnalyzeArgs are argv for the post-record metrics command, e.g.
	// ["metrics","last"] or ["metrics","top","--n","5"].
	AnalyzeArgs []string
	// UseCLI: if true, record via `doctest test` + DOCTEST_METRICS_ROOT;
	// else package runner.RunTest with opts.MetricsRoot.
	UseCLI bool
	// ExtraTestArgs appended to `test` when UseCLI (before dir).
	ExtraTestArgs []string
}

type Response struct {
	// help / skill / analyze process results
	ExitCode int
	Stdout   string
	Stderr   string

	// smoke recording side
	MetricsRoot string
	FixtureDir  string
	RunFiles    []string // new *.jsonl absolute paths after suite
	RecordStdout string
	RecordStderr string
	RecordErr    string
	// analyze after record
	AnalyzeExitCode int
	AnalyzeStdout   string
	AnalyzeStderr   string
	AnalyzeErr      string

	// align_skill_warn
	WarnMessage string
	SkillStdout string
	SkillStderr string
	SkillExit   int
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}
	timeout := req.Timeout
	if timeout == 0 {
		timeout = 90 * time.Second
	}

	switch req.Op {
	case "help":
		if req.Bin == "" {
			return nil, fmt.Errorf("help requires req.Bin")
		}
		code, out, errOut, err := runBin(t, req.Bin, "", nil, req.Args, timeout)
		resp.ExitCode = code
		resp.Stdout = out
		resp.Stderr = errOut
		if err != nil && code == 0 {
			return resp, err
		}
		return resp, nil

	case "align_skill_warn":
		resp.WarnMessage = metrics.FormatDefaultSuiteSlowWarning()
		if req.Bin == "" {
			return nil, fmt.Errorf("align_skill_warn requires req.Bin")
		}
		args := req.Args
		if len(args) == 0 {
			args = []string{"skill", "review-perf", "--show"}
		}
		code, out, errOut, err := runBin(t, req.Bin, "", nil, args, timeout)
		resp.SkillExit = code
		resp.ExitCode = code
		resp.SkillStdout = out
		resp.SkillStderr = errOut
		resp.Stdout = out
		resp.Stderr = errOut
		if err != nil && code == 0 {
			return resp, err
		}
		return resp, nil

	case "smoke":
		root := req.MetricsRoot
		if root == "" {
			root = t.TempDir()
		}
		resp.MetricsRoot = root

		dir := req.FixtureDir
		if dir == "" {
			dir = t.TempDir()
			testtree.WritePassFailTree(t, dir, 1, 0)
			// Seed git origin so project_id is origin-based and matches metrics CLI cwd.
			seedGitOrigin(t, dir, FixtureOrigin)
		}
		resp.FixtureDir = dir

		before := listRunJSONL(t, root)

		if req.UseCLI {
			if req.Bin == "" {
				return nil, fmt.Errorf("smoke UseCLI requires req.Bin")
			}
			args := []string{"test"}
			args = append(args, req.ExtraTestArgs...)
			args = append(args, dir)
			env := []string{"DOCTEST_METRICS_ROOT=" + root}
			code, out, errOut, err := runBin(t, req.Bin, dir, env, args, timeout)
			resp.RecordStdout = out
			resp.RecordStderr = errOut
			if err != nil {
				resp.RecordErr = err.Error()
			}
			if code != 0 && resp.RecordErr == "" {
				resp.RecordErr = fmt.Sprintf("exit %d", code)
			}
		} else {
			var stdout, stderr bytes.Buffer
			opts := core.Options{
				Stdout:      &stdout,
				Stderr:      &stderr,
				RemoveTemp:  true,
				MetricsRoot: root,
			}
			err := runner.RunTest(dir, opts)
			resp.RecordStdout = stdout.String()
			resp.RecordStderr = stderr.String()
			if err != nil {
				resp.RecordErr = err.Error()
			}
		}

		after := listRunJSONL(t, root)
		resp.RunFiles = newPaths(before, after)

		// Analyze with same MetricsRoot; cwd = fixture dir so project_id matches.
		if req.Bin == "" {
			return nil, fmt.Errorf("smoke analyze requires req.Bin")
		}
		analyze := req.AnalyzeArgs
		if len(analyze) == 0 {
			analyze = []string{"metrics", "last"}
		}
		env := []string{"DOCTEST_METRICS_ROOT=" + root}
		code, out, errOut, err := runBin(t, req.Bin, dir, env, analyze, timeout)
		resp.AnalyzeExitCode = code
		resp.ExitCode = code
		resp.AnalyzeStdout = out
		resp.AnalyzeStderr = errOut
		resp.Stdout = out
		resp.Stderr = errOut
		if err != nil {
			resp.AnalyzeErr = err.Error()
			if code == 0 {
				return resp, err
			}
		}
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown Op %q", req.Op)
	}
}

func runBin(t *testing.T, bin, workDir string, extraEnv, args []string, timeout time.Duration) (exitCode int, stdout, stderr string, err error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	if workDir != "" {
		cmd.Dir = workDir
	}
	env := append([]string{}, os.Environ()...)
	env = append(env, extraEnv...)
	cmd.Env = env
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	runErr := cmd.Run()
	stdout = outBuf.String()
	stderr = errBuf.String()
	if runErr == nil {
		return 0, stdout, stderr, nil
	}
	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		return exitErr.ExitCode(), stdout, stderr, nil
	}
	if ctx.Err() != nil {
		return -1, stdout, stderr, ctx.Err()
	}
	return -1, stdout, stderr, runErr
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

func seedGitOrigin(t *testing.T, dir, origin string) {
	t.Helper()
	// git init (ignore if already a repo)
	_ = exec.Command("git", "-C", dir, "init").Run()
	_ = exec.Command("git", "-C", dir, "remote", "remove", "origin").Run()
	if out, err := exec.Command("git", "-C", dir, "remote", "add", "origin", origin).CombinedOutput(); err != nil {
		got, _ := exec.Command("git", "-C", dir, "remote", "get-url", "origin").Output()
		if strings.TrimSpace(string(got)) != origin {
			t.Fatalf("git remote add origin: %v\n%s", err, out)
		}
	}
}

func combinedOut(resp *Response) string {
	return resp.Stdout + "\n" + resp.Stderr
}

func combinedAnalyze(resp *Response) string {
	return resp.AnalyzeStdout + "\n" + resp.AnalyzeStderr
}

// keep metrics import used if leaves only touch helpers
var _ = metrics.SchemaVersion
```
