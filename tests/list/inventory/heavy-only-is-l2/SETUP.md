# Scenario

**Feature**: `label: heavy` without `e2e` is L2; heavy appears in dist

```
Harness -> one leaf label: heavy
  -> list <root>
  -> L2:L3=1:0 (100.0%/0.0%); heavy=1
```

## Steps

1. Write leaf `slow|heavy`.
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
	writeLabeledLeaves(t, root, []string{"slow|heavy"})
	req.FixtureDir = root
	req.Roots = []string{root}
	req.Args = listArgs(nil, root)
	return nil
}
```
