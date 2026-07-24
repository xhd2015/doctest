# Scenario

**Feature**: ASSERT.md path runs labeled leaf explicitly

```
# labeled leaf via ASSERT.md path
doctest test <leaf>/ASSERT.md -> PASS(1/1), no skip summary
```

## Steps

1. Create labeled-only temp tree.
2. Run `doctest test <leaf>/ASSERT.md`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	root := writeLabeledTree(t, false, "ui-automation", "assert path run")
	assertPath := filepath.Join(root, "labeled_leaf", "ASSERT.md")
	req.Args = []string{"test", assertPath}
	return nil
}
```