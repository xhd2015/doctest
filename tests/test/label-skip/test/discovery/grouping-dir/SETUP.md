# Scenario

**Feature**: grouping directory discovery skips labeled child leaves

```
# e2e grouping node, not a leaf
doctest test <tree>/e2e -> PASS(1/1) + skip labeled child
```

## Steps

1. Create tree with `e2e/` grouping containing fast and labeled children.
2. Run `doctest test <tree>/e2e`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	root := writeGroupingLabeledTree(t)
	req.Args = []string{"test", filepath.Join(root, "e2e")}
	return nil
}
```