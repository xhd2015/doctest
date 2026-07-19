# Scenario

**Feature**: flag on does not crash classic gen in P0

```
RunTest(1-leaf, ExperimentRefInsteadOfInline=true) -> still classic gen -> success
```

## Preconditions

- P0 implementer may only set the bool; no ref assembler yet.
- Suite must still complete with exit/run success.

## Steps

1. Run one-leaf pass fixture with field true.
2. Expect no run error.

## Context

- Optional P0 exit criterion: flag-on is parse-safe and run-safe under classic gen.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.ExperimentRefInsteadOfInline = true
	return nil
}
```
