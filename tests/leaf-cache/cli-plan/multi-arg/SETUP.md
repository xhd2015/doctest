# Scenario

**Feature**: multi-arg `doctest test treeA treeB` leaf-cache product path

```
# two explicit DOCTEST roots in one CLI invocation
doctest test tree-a tree-b
  -> union of roots (shared store + skip policy)
  -> PutPass on pass; warm second run Cached > 0
```

## Preconditions

- Parent set binary + isolated leaf-cache env.
- Fixture: all-pass two-tree module (`prepareMultiArgAllPass`).
- Args list two absolute tree roots (not `./...`).

## Steps

1. Build all-pass multi-tree module with tree-a + tree-b.
2. Set multi-arg Args (children may override Args2/Args3).
3. Assert warm / disable outcomes under multi-arg.

## Context

- Sibling of `workspace/` (which uses `<mod>/...`). Same fixture shape, different
  CLI shape — proves multi-arg is not a second engine with different physics.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.Op = "runtime_multi"
	prepareMultiArgAllPass(t, req)
	args := multiArgTwoTrees(req)
	req.Args = args
	req.Args2 = append([]string(nil), args...)
	return nil
}
```
