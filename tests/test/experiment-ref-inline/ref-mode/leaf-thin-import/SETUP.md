# Scenario

**Feature**: leaf packages are thin imports of root (no inlined types/helpers)

```
RunTest(2-leaf, GenDir=tmp)
  -> leaf packages under a/ b/ lack marker def and type Request
  -> each imports a non-stdlib package
```

## Preconditions

- Default hierarchical unified path.
- Leaves may be non-`_test` RunTestLeaf packages.

## Steps

1. `Op=ref_gen`.
2. Assert thin leaf properties.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "ref_gen"
	return nil
}
```
