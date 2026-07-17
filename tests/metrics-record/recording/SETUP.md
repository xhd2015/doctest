# Scenario

**Feature**: suite-level metrics JSONL under injectable MetricsRoot

```
# metrics on
RunTest(dir, Options{MetricsRoot, NoMetrics:false})
  -> $MetricsRoot/doctest/metrics/<project>/runs/*.jsonl
  -> run_start ... run_end [leaf_*]

# metrics off
RunTest(..., NoMetrics:true) -> no new *.jsonl under MetricsRoot
```

## Preconditions

- Fixture is a tiny 1-leaf pass tree (`testtree.WritePassFailTree`) unless a leaf overrides `req.Dir`.
- `MetricsRoot` is always a fresh `t.TempDir()` for isolation.
- Package path preferred: `runner.RunTest` with `core.Options.MetricsRoot` / `NoMetrics`.
  CLI path optional via `UseCLI` + `DOCTEST_METRICS_ROOT` + `req.Bin` (set only if a leaf enables CLI).

## Steps

1. Leaf sets `req.Op = "record_run"` and metrics on/off.
2. Run suite against fixture.
3. Diff `*.jsonl` under MetricsRoot; decode events when present.

## Context

- Project id may be nogit_* when the temp fixture has no origin — still valid P1 layout.
- Leaves should not require wall time near 3 minutes.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.Op = "record_run"
	if req.MetricsRoot == "" {
		req.MetricsRoot = t.TempDir()
	}
	req.UseCLI = false
	return nil
}
```
