# Scenario

**Feature**: flag off keeps classic per-leaf inline generation

```
RunTest(2-leaf, ExperimentRefInsteadOfInline=false, GenDir=tmp)
  -> classic AssembleTestSource
  -> each leaf package inlines root helper
```

## Preconditions

- Field explicitly false (default path).

## Steps

1. Set `ExperimentRefInsteadOfInline=false`.
2. Assert classic layout / pass.

## Context

- Regression guard for hard product rule: without the flag, no package DAG.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.ExperimentRefInsteadOfInline = false
	return nil
}
```
