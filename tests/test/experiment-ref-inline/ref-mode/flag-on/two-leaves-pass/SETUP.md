# Scenario

**Feature**: simple 2-leaf tree passes under ref mode

```
doctest test --experiment-ref-instead-of-inline --gen-dir tmp fixture/{a,b}
  -> both leaves pass (exit/run ok)
```

## Preconditions

- Default two-leaf marker fixture.

## Steps

1. Run with flag on.
2. Expect no `RunErr`.

## Context

- Exit criterion: both leaves pass under the flag (does not inspect layout).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.ExperimentRefInsteadOfInline = true
	req.Op = "ref_gen"
	return nil
}
```

