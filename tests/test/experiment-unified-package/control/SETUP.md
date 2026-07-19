# Scenario

**Feature**: unified flag off keeps classic multi-leaf generation

```
RunTest(2-leaf, ExperimentUnifiedPackagePerDoctestTree=false, GenDir=tmp)
  -> classic per-leaf *_test.go packages
```

## Preconditions

- Uses same fixture factory as unified-mode for an apples-to-apples control.
- Unified Options field explicitly false.

## Steps

1. Set `Op=run_gen` and unified false.
2. Assert classic leaf `*_test.go` layout.

## Context

- Hard product rule: without the unified flag, behavior is unchanged.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "run_gen"
	req.ExperimentUnifiedPackagePerDoctestTree = false
	return nil
}
```
