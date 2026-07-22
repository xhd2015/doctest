# Scenario

**Feature**: same relative leaf path across workspace trees must not false-skip

```
# tree-a/leaf pass, tree-b/leaf fail (identical relpath "leaf")
run1: doctest test tree-a     -> PutPass only tree-a identity
run2: doctest test <mod>/...  -> tree-a may Cached; tree-b MUST execute (fail)
```

## Preconditions

- Same-relpath pass/fail fixture (`prepareWorkspaceSameRelpathPassFail`).
- Shared leaf-cache store env across both runs.
- Bare relative skip tokens would false-skip tree-b → exit 0 (bug).

## Steps

1. Build fixture.
2. Args = single-tree `test tree-a` (warm only tree-a).
3. Args2 = workspace `test <mod>/...`.
4. Assert run2 fails and Cached == 1 (only tree-a).

## Context

- Complements `polish/isolation/same-relpath-two-trees` (separate single-tree
  runs) with a **shared workspace suite** skip-env collision case.
- Implementer must use tree-qualified skip identities (FormatLeafIdentity or
  equivalent) so tree-b's body still runs.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	prepareWorkspaceSameRelpathPassFail(t, req)
	// Warm only tree-a (single-tree path — already GREEN for leaf-cache).
	req.Args = []string{"test", req.TreeRoot}
	// Workspace multi-tree: must not treat bare "leaf" as skip for tree-b.
	req.Args2 = []string{"test", mustWorkspacePattern(req.WorkDir)}
	return nil
}
```
