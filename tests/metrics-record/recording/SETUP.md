# Scenario

**Feature**: suite metrics recording under injectable MetricsRoot

```
RunTest(dir, Options{MetricsRoot, MetricsOn:true})
  -> run JSONL with start/end/leaf events

RunTest(..., MetricsOn:false) -> no new *.jsonl under MetricsRoot
```

## Preconditions

- Package path preferred: `runner.RunTest` with `core.Options.MetricsRoot` / `MetricsOn`.

## Steps

1. Leaf sets MetricsOn and optional fixture.
2. Assert JSONL presence or absence.

## Context

- MetricsRoot is always a temp dir in these leaves.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "record_run"
	return nil
}
```
