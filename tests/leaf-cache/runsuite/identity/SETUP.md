# Scenario

**Feature**: tree-qualified leaf identities for multi-tree skip/fail maps

```
FormatLeafIdentity(treeA, "leaf") -> idA
FormatLeafIdentity(treeB, "leaf") -> idB
# idA != idB (tree-qualified)
# FormatLeafIdentity stable for identical inputs
```

## Preconditions

- Twin trees from `prepareTwinTrees` (same relative leaf name `leaf`).
- Op=`format_identity` for descendant leaves unless overridden.

## Steps

1. Build twin trees.
2. Format identities for each tree's relative leaf path.
3. Assert non-empty, stable, and cross-tree distinct.

## Context

- Complements `key/tree-identity` (store **keys** include TreeRoot) with the
  in-memory **identity** token used by prepare/record maps.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "format_identity"
	prepareTwinTrees(t, req)
	return nil
}
```
