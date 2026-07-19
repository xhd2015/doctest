# Scenario

**Feature**: unified gen layout under GenDir

```
RunTest(2-leaf, unified=true, GenDir=tmp)
  -> walk GenDir
  -> __droot, __registry, __allleaves, suite present
  -> leaf packages are non-_test with RunTestLeaf
  -> no leaf a|b *_test.go
  -> suite imports only __registry + __allleaves
```

## Preconditions

- Same default fixture and unified Options as siblings.

## Steps

1. Run with unified flag on.
2. Assert directory/file shape and suite imports from `Response` layout fields.

## Context

- Layout names are locked: `__droot`, `__registry`, `__allleaves`, `suite`.
- Leaf entrypoint name locked: `RunTestLeaf`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "run_gen"
	req.ExperimentUnifiedPackagePerDoctestTree = true
	return nil
}
```
