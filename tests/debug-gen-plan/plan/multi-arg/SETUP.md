# Scenario

**Feature**: multi-arg plan prints per-arg package trees then a merged view

```
doctest test tree-a tree-b --gen-dir G
  -> gen-plan: arg[1/2] tree-a   # package subtree only
  -> gen-plan: arg[2/2] tree-b
  -> gen-plan: merged
       go.mod / go.sum / manifest / tidy-done
       <both trees>
       [__workspace if hub written]
```

## Preconditions

- Two DOCTEST roots under one module; explicit multi-arg CLI.

## Steps

1. prepareMultiArgTwoTrees; Args = test tree-a tree-b with gen-dir.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	prepareMultiArgTwoTrees(t, req)
	req.Args = baseTestArgs(req, "--no-color", "tree-a", "tree-b")
	return nil
}
```
