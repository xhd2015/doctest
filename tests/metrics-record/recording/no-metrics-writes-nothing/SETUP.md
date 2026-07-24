# Scenario

**Feature**: default MetricsOn=false writes no run file under MetricsRoot

```
RunTest(fixture, MetricsOn=false, MetricsRoot=tmp) -> no new *.jsonl
```

## Preconditions

- Empty MetricsRoot before run.

## Steps

1. Run suite with MetricsOn false (default).
2. Assert no new run JSONL files.

## Context

- Contrast with enabled-writes leaves that set MetricsOn=true.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.MetricsOn = false
	return nil
}
```
