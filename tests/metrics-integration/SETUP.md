# Scenario

**Feature**: end-to-end metrics polish — record, analyze, help, skill alignment

```
# smoke: tiny suite records under injectable root; metrics CLI reads it
1-leaf pass tree
  -> RunTest / doctest test (MetricsRoot | DOCTEST_METRICS_ROOT)
  -> runs/*.jsonl
  -> doctest metrics last|top -> run evidence

# help surfaces
doctest --help -> metrics
doctest test --help -> --no-metrics

# skill ↔ WARNING
FormatDefaultSuiteSlowWarning phrases ⊆ skill review-perf --show
```

## Preconditions

- Module root is `DOCTEST_ROOT/../..`.
- P1–P4 wired: recording, metrics CLI, review-perf skill, warn helpers.
- Leaves never sleep for the 3-minute budget; smoke suites are one fast leaf.
- MetricsRoot is always `t.TempDir()` so the user cache is untouched.
- Smoke uses the **fixture directory as metrics CLI cwd** so project_id matches
  the recorder’s `projectIDForDir(fixture)`.

## Steps

1. Build (or reuse) the doctest binary via `testbin.Ensure` when Bin is needed.
2. Leaf Setup sets `req.Op` and branch fields (AnalyzeArgs, Args).
3. Root `Run` dispatches smoke / help / align_skill_warn.

## Context

- Prefer package `runner.RunTest` for recording (faster, no nested CLI rebuild of
  the fixture suite beyond what RunTest already does). Leaves may set `UseCLI`
  if they specifically need the env-injection CLI path.
- Help and skill leaves only need the binary.
- Parallel-safe: per-leaf temp MetricsRoot and fixture dirs.

```go
import (
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
)

func Setup(t *testing.T, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 90 * time.Second
	}
	// Smoke analyze + help + skill need the CLI binary.
	if req.Bin == "" {
		req.Bin = testbin.Ensure(t, filepath.Join(DOCTEST_ROOT, "..", ".."))
	}
	return nil
}
```
