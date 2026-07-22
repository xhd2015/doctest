# Scenario

**Feature**: workspace `doctest test -v ./...` prints planned trees+tests before go test

```
module with tree_a (1 leaf) + tree_b (1 leaf)
doctest test -v --no-color ./...
  -> doctest: workspace (2 trees, 2 tests)   # or hub label
  -> cd … && go test …
```

## Preconditions

- WorkDir = module root; args use `./...`.
- Classic TDD: **RED** until verbose workspace path prints the planned line
  (quiet already does).

## Steps

1. Create two-tree workspace module.
2. Run `doctest test -v --no-color ./...` from the module root.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	mod := createWorkspaceTwoTrees(t)
	req.WorkDir = mod
	req.Args = []string{"test", "-v", "--no-color", "./..."}
	return nil
}
```
