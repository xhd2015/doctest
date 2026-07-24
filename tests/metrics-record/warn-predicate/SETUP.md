# Scenario

**Feature**: pure predicate for default-suite slow WARNING

```
# decide without I/O or sleep
ShouldWarnDefaultSuiteSlow(default_suite, total, elapsed, threshold) -> bool

# true only when default_suite && total > 0 && elapsed > threshold
```

## Preconditions

- Leaves set synthetic elapsed; never wait wall time.
- Threshold defaults to `metrics.DefaultSuiteWarnThreshold` (3m).

## Steps

1. Set `req.Op = "should_warn"` and predicate inputs.
2. Assert boolean outcome.

## Context

- Not default suite includes LabelAll and non-empty LabelExprs (modeled here as
  `DefaultSuite=false`).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "should_warn"
	return nil
}
```
