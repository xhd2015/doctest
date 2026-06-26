# Scenario

**Feature**: concrete leaf path executes labeled test

```
# labeled leaf dir
doctest test <tree>/labeled_leaf -> PASS(1/1), no skip summary
```

## Steps

1. Create labeled-only temp tree.
2. Run `doctest test <tree>/labeled_leaf`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	root := writeLabeledTree(t, false, "ui-automation", "heavy ui test")
	req.Args = []string{"test", filepath.Join(root, "labeled_leaf")}
	return nil
}
```