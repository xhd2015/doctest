# Scenario

**Feature**: editing a local package imported by the spine invalidates the leaf pass

```
run1: test leaf(import helper) -> PutPass
run2: warm -> Cached > 0
mutate pkg/helper/helper.go
run3: 0 Cached (must re-run)
```

## Preconditions

- Fixture module `example.com/app` with `pkg/helper` imported by leaf ASSERT.
- Mutation `polish_edit_local_dep` after run2 (changes Answer from 42 to 99).

## Steps

1. prepareLocalDepPassFixture (assert expects `helper.Answer()==42`).
2. Three default test runs; MutateAfterRun=2 changes Answer to 99.
3. Assert run2 Cached>0; run3 Cached==0 (re-execution). Run3 may fail the value check — that still proves the leaf was not skipped.

## Context

- A false warm hit on run3 (Cached>0 after dep edit) is the bug this leaf guards against.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	dir := prepareLocalDepPassFixture(t, req)
	req.Args = []string{"test", dir}
	req.Args2 = []string{"test", dir}
	req.Args3 = []string{"test", dir}
	req.Mutation = "polish_edit_local_dep"
	req.MutateAfterRun = 2
	return nil
}
```
