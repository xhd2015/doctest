# Scenario

**Feature**: `--no-metrics` / NoMetrics prevents any run file under MetricsRoot

```
# opt-out recording
RunTest(fixture, NoMetrics=true, MetricsRoot=tmp) -> no new *.jsonl
```

## Preconditions

- Same 1-leaf pass fixture as enabled cases.

## Steps

1. Run suite with NoMetrics true.
2. Assert zero new JSONL files under MetricsRoot.

## Context

- Suite may still pass tests; only metrics I/O is suppressed.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.NoMetrics = true
	return nil
}
```
