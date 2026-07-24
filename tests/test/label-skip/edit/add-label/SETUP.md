# Scenario

**Feature**: edit creates frontmatter with label and explanation

```
# unlabeled leaf
doctest edit <leaf> --add-label --add-explanation -> ASSERT.md frontmatter
```

## Steps

1. Create unlabeled temp tree.
2. Run `doctest edit <leaf> --add-label ui-automation --add-explanation "AX test"`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	root := writeUnlabeledTree(t)
	leaf := filepath.Join(root, "plain_leaf")
	req.Args = []string{"edit", leaf, "--add-label", "ui-automation", "--add-explanation", "AX test"}
	return nil
}
```