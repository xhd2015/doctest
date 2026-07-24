# Scenario

**Feature**: grouping SubDir discovery skips labeled child

```
FilterBySubDir(e2e) + discovery partition → run fast_child; skip labeled_child
```

## Steps

1. Grouping fixture; SubDir = e2e.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.TreeRoot = writeGroupingLabeledTree(t)
	req.SubDir = filepath.Join(req.TreeRoot, "e2e")
	return nil
}
```
