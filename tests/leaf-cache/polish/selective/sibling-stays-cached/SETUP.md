# Scenario

**Feature**: editing one leaf ASSERT changes only that leaf's key; sibling key stable

```
keys0 = ComputeLeafKey(leaf_a), ComputeLeafKey(leaf_b)
mutate leaf_a ASSERT Go
keys1 = ...
# leaf_a key changes; leaf_b key unchanged
```

## Preconditions

- Fixture with `leaf_a` and `leaf_b` (`prepareTwoSiblingPassLeaves`).
- Mutation `polish_edit_leaf_a`.

## Steps

1. Prepare two-sibling fixture.
2. Op=`two_sibling_keys`.

## Context

- Library proof that keys are leaf-scoped (product Cached==1 is the L3 consequence under runtime warm + multi-leaf).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "two_sibling_keys"
	req.FixtureDir = prepareTwoSiblingPassLeaves(t)
	req.TreeRoot = req.FixtureDir
	req.ModuleRoot = req.FixtureDir
	req.Mutation = "polish_edit_leaf_a"
	return nil
}
```
