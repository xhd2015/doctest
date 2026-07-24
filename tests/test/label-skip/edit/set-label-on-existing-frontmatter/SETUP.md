# Scenario

**Feature**: edit appends a second label to existing frontmatter

```
# one label already present
doctest edit --add-label manual -> comma-separated labels
```

## Steps

1. Create tree with label `ui-automation`.
2. Run `doctest edit <leaf> --add-label manual`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	root := writeLabeledTree(t, false, "ui-automation", "first label")
	leaf := filepath.Join(root, "labeled_leaf")
	req.Args = []string{"edit", leaf, "--add-label", "manual"}
	return nil
}
```