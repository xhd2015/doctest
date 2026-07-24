# Scenario

**Feature**: default gen layout is hierarchical unified packages

```
RunTest(2-leaf, GenDir=tmp)
  -> __droot, __registry, __allleaves, suite (runall.go + suite_test.go)
  -> leaf non-test RunTestLeaf; no leaf *_test.go
```

## Preconditions

- Default generation (no experiment flags).
- GenDir kept for walk.

## Steps

1. `Op=run_gen`.
2. Assert layout fields on Response.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "run_gen"
	return nil
}
```
