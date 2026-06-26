# Scenario

**Feature**: edit accepts ASSERT.md path same as leaf directory

```
# ASSERT.md path instead of leaf dir
doctest edit <leaf>/ASSERT.md --add-label -> frontmatter updated
```

## Steps

1. Create unlabeled temp tree.
2. Run `doctest edit <leaf>/ASSERT.md --add-label ui-automation`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	root := writeUnlabeledTree(t)
	assertPath := filepath.Join(root, "plain_leaf", "ASSERT.md")
	req.Args = []string{"edit", assertPath, "--add-label", "ui-automation"}
	return nil
}
```