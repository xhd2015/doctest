# Scenario

**Feature**: edit warns when adding an already-present label

```
# labeled leaf, duplicate --add-label
doctest edit --add-label existing -> stderr warning, ASSERT.md unchanged
```

## Steps

1. Create labeled temp tree and record ASSERT.md content.
2. Run `doctest edit <leaf> --add-label ui-automation`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	root := writeLabeledTree(t, false, "ui-automation", "heavy ui test")
	leaf := filepath.Join(root, "labeled_leaf")
	req.Args = []string{"edit", leaf, "--add-label", "ui-automation"}
	return nil
}
```