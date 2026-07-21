# Scenario

**Feature**: edit package only leaf-d uses → peer leaves stay Cached

```
run1: -count=1 -> 3 Run, 0 Cached (seed store)
edit alone/d Version
run2: no -count -> 3 Run, 2 Cached (leaf-ab-1 + leaf-ab-2 warm)
```

## Preconditions

- `preparePartialPackageDepsFixture` layout.
- Mutation `polish_edit_alone_d` after run1.

## Steps

1. Build fixture.
2. Args = `test -count=1 <tree>`; Args2 = `test <tree>` (no count).
3. MutateAfterRun=1.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	dir := preparePartialPackageDepsFixture(t, req)
	req.Args = []string{"test", "-count=1", dir}
	req.Args2 = []string{"test", dir}
	req.Mutation = "polish_edit_alone_d"
	req.MutateAfterRun = 1
	return nil
}
```
