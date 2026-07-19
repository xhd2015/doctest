# Scenario

**Feature**: flag off classic multi-package leaf tests

```
doctest test --gen-dir tmp fixture/{a,b}
  -> leaf a and b each have *_test.go
  -> run succeeds
  -> go test package args are multi-leaf (not suite-only)
```

## Preconditions

- Default two-leaf marker fixture; unified false from ancestor.

## Steps

1. Run with unified off.
2. Expect classic leaf `*_test.go` and multi-package shape.

## Context

- Control leaf: proves unified is opt-in only.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Explicit classic control: unified off, run_gen for layout + pass.
	req.Op = "run_gen"
	req.ExperimentUnifiedPackagePerDoctestTree = false
	return nil
}
```
