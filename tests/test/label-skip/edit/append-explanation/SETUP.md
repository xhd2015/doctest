# Scenario

**Feature**: edit appends explanation with semicolon separator

```
# existing explanation
doctest edit --add-explanation second -> "first; second"
```

## Steps

1. Create tree with label and explanation `first`.
2. Run `doctest edit <leaf> --add-explanation second`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	root := writeLabeledTree(t, false, "ui-automation", "first")
	leaf := filepath.Join(root, "labeled_leaf")
	req.Args = []string{"edit", leaf, "--add-explanation", "second"}
	return nil
}
```