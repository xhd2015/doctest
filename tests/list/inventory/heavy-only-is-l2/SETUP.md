# Scenario

**Feature**: `label: slow` without `e2e` is L2; slow appears in dist

```
Harness -> one leaf label: slow
  -> list <root>
  -> L2:L3=1:0 (100.0%/0.0%); slow=1
```

## Steps

1. Write leaf `slow|slow`.
2. Args = `list <root>`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	root := filepath.Join(t.TempDir(), "tree")
	writeLabeledLeaves(t, root, []string{"slow|slow"})
	req.FixtureDir = root
	req.Roots = []string{root}
	req.Args = listArgs(nil, root)
	return nil
}
```
