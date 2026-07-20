# Scenario

**Feature**: pass cache does not cross doctest tree boundaries

```
# same relative leaf path "leaf", different abs TreeRoot
warm treeA -> GetPass hits for treeA only
first run treeB -> 0 Cached (no collision)
```

## Preconditions

- Twin trees from prepareTwinTrees; shared DOCTEST_LEAF_CACHE via isolateRuntimeEnv.

## Steps

1. Run1+Run2: test treeA (populate + warm).
2. Run3: test treeB (must be cold).

## Context

- Complements key/tree-identity unit leaf with end-to-end suite behavior.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	return nil
}
```
