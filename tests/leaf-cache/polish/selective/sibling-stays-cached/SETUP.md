# Scenario

**Feature**: editing one leaf ASSERT re-runs only that leaf; sibling stays Cached

```
run1: test 2-leaf tree -> store both
run2: warm -> Cached == 2
mutate leaf_a ASSERT
run3: Cached == 1 (leaf_b only)
```

## Preconditions

- Fixture with `leaf_a` and `leaf_b`.
- Mutation `polish_edit_leaf_a` after run2.

## Steps

1. prepareTwoSiblingPassLeaves.
2. Args/Args2/Args3 = `test <fixture>`; MutateAfterRun=2.

## Context

- Proves keys are leaf-scoped, not tree-scoped wipe on any edit.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	req.FixtureDir = prepareTwoSiblingPassLeaves(t)
	req.Args = []string{"test", req.FixtureDir}
	req.Args2 = []string{"test", req.FixtureDir}
	req.Args3 = []string{"test", req.FixtureDir}
	req.Mutation = "polish_edit_leaf_a"
	req.MutateAfterRun = 2
	return nil
}
```
