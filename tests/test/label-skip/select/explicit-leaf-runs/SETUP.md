# Scenario

**Feature**: ExplicitLeaf runs labeled leaf without LabelAll

```
SubDir=labeled_leaf + ExplicitLeaf → run labeled leaf
```

## Steps

1. Labeled-only tree; SubDir to labeled leaf; ExplicitLeaf true.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeLabeledTree(t, false, "ui-automation", "heavy ui test")
	req.SubDir = filepath.Join(req.TreeRoot, "labeled_leaf")
	req.ExplicitLeaf = true
	return nil
}
```
