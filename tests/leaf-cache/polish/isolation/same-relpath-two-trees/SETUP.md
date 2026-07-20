# Scenario

**Feature**: treeB does not inherit warm skips from treeA at the same relative path

```
run1: test treeA -> pass store
run2: test treeA -> Cached > 0
run3: test treeB -> 0 Cached (first time for treeB)
```

## Preconditions

- Content-identical twins; shared leaf-cache store env.
- Relative leaf path is `leaf` in both trees.

## Steps

1. prepareTwinTrees.
2. Args/Args2 target FixtureDir (A); Args3 targets FixtureB.

## Context

- If keys omit tree identity, run3 may incorrectly show Cached > 0 (bug).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	prepareTwinTrees(t, req)
	req.Args = []string{"test", req.FixtureDir}
	req.Args2 = []string{"test", req.FixtureDir}
	req.Args3 = []string{"test", req.FixtureB}
	return nil
}
```
