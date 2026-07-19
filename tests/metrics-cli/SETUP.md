# Scenario

**Feature**: analyze recorded metrics via `doctest metrics` against fixture JSONL

```
# build CLI once
go build ./cmd/doctest -> req.Bin

# inject metrics root; resolve project from cwd origin
DOCTEST_METRICS_ROOT=<temp>
cwd (git origin) -> project_id
  -> $root/doctest/metrics/<project_id>/runs/*.jsonl
  -> doctest metrics <subcmd> [flags]
```

## Preconditions

- Module root is `DOCTEST_ROOT/../..`.
- P1 layout and P2 event shapes are fixed; this tree only reads fixtures.
- Env `DOCTEST_METRICS_ROOT` overrides the metrics cache root for the CLI
  (same as P2 recording tests).
- Prune retention under test: keep newest **30** run files per project.
- Leaves never invoke long `doctest test` suites to populate metrics; they
  write JSONL fixtures under a temp MetricsRoot.

## Steps

1. Build (or reuse) the doctest binary via `testbin.Ensure`.
2. Descendants create WorkDir + MetricsRoot, seed `runs/*.jsonl`, set `Args`.
3. `Run` executes `req.Bin` with `DOCTEST_METRICS_ROOT` and captures I/O.

## Context

- Shared helpers live in this root `SETUP.md` Go block (fixture writers, git cwd).
- `Request` / `Response` / `Run` are defined only in `DOCTEST.md`.
- Parallel-safe: each leaf uses `t.TempDir()` for MetricsRoot and WorkDir.

```go
import (
"github.com/xhd2015/doctest/session"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/metrics"
	"github.com/xhd2015/doctest/libdoc/testbin"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 45 * time.Second
	}
	if req.Bin == "" {
		req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	}
	return nil
}

// ensureFixtureProject creates a temp git cwd with FixtureOrigin and sets
// req.WorkDir, req.ProjectID, and (if empty) req.MetricsRoot.
func ensureFixtureProject(t *testing.T, req *Request) {
	t.Helper()
	if req.MetricsRoot == "" {
		req.MetricsRoot = t.TempDir()
	}
	if req.WorkDir == "" {
		req.WorkDir = t.TempDir()
	}
	if req.ProjectID == "" {
		req.ProjectID = metrics.ProjectIDFromOrigin(FixtureOrigin)
	}
	// git init + origin (ignore errors if already initialized)
	runQuiet(t, req.WorkDir, "git", "init")
	runQuiet(t, req.WorkDir, "git", "remote", "remove", "origin")
	if out, err := exec.Command("git", "-C", req.WorkDir, "remote", "add", "origin", FixtureOrigin).CombinedOutput(); err != nil {
		// if origin already correct, ok; else fail
		got, _ := exec.Command("git", "-C", req.WorkDir, "remote", "get-url", "origin").Output()
		if strings.TrimSpace(string(got)) != FixtureOrigin {
			t.Fatalf("git remote add origin: %v\n%s", err, out)
		}
	}
}

func runQuiet(t *testing.T, dir string, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	_ = cmd.Run()
}

func projectRunsDir(req *Request) string {
	return filepath.Join(req.MetricsRoot, "doctest", "metrics", req.ProjectID, "runs")
}

func projectMetricsDir(req *Request) string {
	return filepath.Join(req.MetricsRoot, "doctest", "metrics", req.ProjectID)
}

// writeRunFile writes newline-delimited JSON objects to runs/<name>.
// name should be a full basename including .jsonl.
func writeRunFile(t *testing.T, req *Request, name string, events []map[string]any) string {
	t.Helper()
	dir := projectRunsDir(req)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir runs: %v", err)
	}
	path := filepath.Join(dir, name)
	var b strings.Builder
	for _, ev := range events {
		if _, ok := ev["schema_version"]; !ok {
			ev["schema_version"] = metrics.SchemaVersion
		}
		data, err := json.Marshal(ev)
		if err != nil {
			t.Fatalf("marshal event: %v", err)
		}
		b.Write(data)
		b.WriteByte('\n')
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write run file: %v", err)
	}
	return path
}

func runStem(name string) string {
	return strings.TrimSuffix(name, ".jsonl")
}

// fixtureRunDefault builds a default-suite run with mixed labeled/unlabeled leaves.
// elapsed ranking (desc): slow-leaf, labeled-leaf, mid-leaf, fast-leaf.
func fixtureRunDefault(runID string) []map[string]any {
	return []map[string]any{
		{
			"type":       "run_start",
			"run_id":     runID,
			"ts":         "2026-07-16T09:00:00Z",
			"project_id": metrics.ProjectIDFromOrigin(FixtureOrigin),
			"cwd":        "/work/fixture",
			"argv":       []any{"doctest", "test", "./tests"},
			"mode": map[string]any{
				"default_suite": true,
				"label_all":     false,
				"label_exprs":   []any{},
			},
		},
		{
			"type":     "leaf_start",
			"path":     "group/slow-leaf",
			"root":     "tests/fixture",
			"ts_start": "2026-07-16T09:00:01Z",
			"labels":   []any{},
		},
		{
			"type":       "leaf_end",
			"path":       "group/slow-leaf",
			"ts_end":     "2026-07-16T09:00:06Z",
			"elapsed_ns": int64(5_000_000_000),
			"result":     "pass",
			"cached":     false,
		},
		{
			"type":     "leaf_start",
			"path":     "group/mid-leaf",
			"root":     "tests/fixture",
			"ts_start": "2026-07-16T09:00:06Z",
			"labels":   []any{},
		},
		{
			"type":       "leaf_end",
			"path":       "group/mid-leaf",
			"ts_end":     "2026-07-16T09:00:08Z",
			"elapsed_ns": int64(2_000_000_000),
			"result":     "pass",
			"cached":     false,
		},
		{
			"type":     "leaf_start",
			"path":     "group/fast-leaf",
			"root":     "tests/fixture",
			"ts_start": "2026-07-16T09:00:08Z",
			"labels":   []any{},
		},
		{
			"type":       "leaf_end",
			"path":       "group/fast-leaf",
			"ts_end":     "2026-07-16T09:00:08.1Z",
			"elapsed_ns": int64(100_000_000),
			"result":     "pass",
			"cached":     false,
		},
		{
			"type":     "leaf_start",
			"path":     "group/labeled-leaf",
			"root":     "tests/fixture",
			"ts_start": "2026-07-16T09:00:08.1Z",
			"labels":   []any{"slow"},
		},
		{
			"type":       "leaf_end",
			"path":       "group/labeled-leaf",
			"ts_end":     "2026-07-16T09:00:11Z",
			"elapsed_ns": int64(3_000_000_000),
			"result":     "pass",
			"cached":     false,
		},
		{
			"type":      "run_end",
			"wall_ns":   int64(11_000_000_000),
			"passed":    4,
			"total":     4,
			"skipped":   0,
			"exit_ok":   true,
			"warnings":  []any{},
		},
	}
}

// fixtureRunLabeledSuite is a non-default (label-filtered) suite with one slow leaf.
func fixtureRunLabeledSuite(runID string) []map[string]any {
	return []map[string]any{
		{
			"type":   "run_start",
			"run_id": runID,
			"ts":     "2026-07-15T12:00:00Z",
			"mode": map[string]any{
				"default_suite": false,
				"label_all":     false,
				"label_exprs":   []any{"slow"},
			},
		},
		{
			"type":     "leaf_start",
			"path":     "only/labeled-suite-leaf",
			"ts_start": "2026-07-15T12:00:01Z",
			"labels":   []any{"slow"},
		},
		{
			"type":       "leaf_end",
			"path":       "only/labeled-suite-leaf",
			"elapsed_ns": int64(9_000_000_000),
			"result":     "pass",
		},
		{
			"type":     "run_end",
			"wall_ns":  int64(9_500_000_000),
			"passed":   1,
			"total":    1,
			"skipped":  0,
			"exit_ok":  true,
			"warnings": []any{},
		},
	}
}

// fixtureRunOlder is a smaller older default-suite run (for last/show multi-file).
func fixtureRunOlder(runID string) []map[string]any {
	return []map[string]any{
		{
			"type":   "run_start",
			"run_id": runID,
			"ts":     "2026-07-01T10:00:00Z",
			"mode": map[string]any{
				"default_suite": true,
			},
		},
		{
			"type":     "leaf_start",
			"path":     "old/only-leaf",
			"labels":   []any{},
			"ts_start": "2026-07-01T10:00:01Z",
		},
		{
			"type":       "leaf_end",
			"path":       "old/only-leaf",
			"elapsed_ns": int64(500_000_000),
			"result":     "pass",
		},
		{
			"type":     "run_end",
			"wall_ns":  int64(600_000_000),
			"passed":   1,
			"total":    1,
			"skipped":  0,
			"exit_ok":  true,
			"warnings": []any{},
		},
	}
}

func mustContain(t *testing.T, hay, needle, label string) {
	t.Helper()
	if !strings.Contains(hay, needle) {
		t.Fatalf("%s: missing %q in:\n%s", label, needle, hay)
	}
}

func combinedOut(resp *Response) string {
	return resp.Stdout + "\n" + resp.Stderr
}

func isValidJSON(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	var v any
	return json.Unmarshal([]byte(s), &v) == nil
}

// silence unused import if a leaf only uses helpers through side effects
var _ = fmt.Sprintf
```
