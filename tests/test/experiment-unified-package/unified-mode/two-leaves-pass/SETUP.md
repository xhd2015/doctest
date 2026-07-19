# Scenario

**Feature**: simple 2-leaf tree passes under unified package mode

```
doctest test --experiment-unified-package-per-doctest-tree --gen-dir tmp fixture/{a,b}
  -> both leaves pass (exit/run ok)
```

## Preconditions

- Default two-leaf marker fixture.

## Steps

1. Run with unified flag on (via Options).
2. Expect no `RunErr`.

## Context

- Exit criterion: both leaves pass under the flag (does not inspect layout).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "run_gen"
	req.ExperimentUnifiedPackagePerDoctestTree = true
	return nil
}
```
