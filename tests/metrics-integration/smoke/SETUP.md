# Scenario

**Feature**: tiny suite under MetricsRoot produces JSONL that metrics CLI can read

```
# record
1-leaf pass tree + MetricsRoot
  -> RunTest (metrics on)
  -> ≥1 *.jsonl under MetricsRoot

# analyze
cwd=fixture + DOCTEST_METRICS_ROOT
  -> doctest metrics last|top -> exit 0 + evidence of the run
```

## Preconditions

- Metrics must be opt-in (`MetricsOn=true` / `--metrics-on`).
- Fixture is a minimal pass tree (`testtree.WritePassFailTree` 1 pass) with a
  seeded git `origin` (`FixtureOrigin`) so record and analyze share one project_id.
- Analyze cwd equals fixture dir for project_id alignment.

## Steps

1. Leaf sets `Op=smoke` and `AnalyzeArgs`.
2. Run records then analyzes.
3. Assert JSONL present and analyze output useful.

## Context

- Package `RunTest` path is default (`UseCLI=false`).
- Leaf path from the fixture is typically under `a_pass_0` (WritePassFailTree).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "smoke"
	req.UseCLI = false
	return nil
}
```
