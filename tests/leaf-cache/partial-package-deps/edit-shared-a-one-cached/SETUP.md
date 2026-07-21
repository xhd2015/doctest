# Scenario

**Feature**: edit shared package used by two leaves → only leaf-d stays Cached

```
run1: -count=1 -> 0 Cached
edit shared/a Version
run2: no -count -> 1 Cached (leaf-d only)
```

## Preconditions

- Same partial-package fixture as sibling leaf.
- Mutation `polish_edit_shared_a` after run1.

## Steps

1. Build fixture.
2. Run1 `-count=1`; Run2 without count; mutate shared/a between.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	dir := preparePartialPackageDepsFixture(t, req)
	req.Args = []string{"test", "-count=1", dir}
	req.Args2 = []string{"test", dir}
	req.Mutation = "polish_edit_shared_a"
	req.MutateAfterRun = 1
	return nil
}
```
