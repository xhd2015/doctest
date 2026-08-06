# Scenario

**Feature**: multi-label `e2e, slow` counts as L3; both labels +1 in dist

```
Harness -> one leaf label: e2e, slow
  -> list <root>
  -> L2:L3=0:1 (0.0%/100.0%); e2e=1 slow=1
```

## Steps

1. Write leaf `both|e2e, slow`.
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
	writeLabeledLeaves(t, root, []string{"both|e2e, slow"})
	req.FixtureDir = root
	req.Roots = []string{root}
	req.Args = listArgs(nil, root)
	return nil
}
```
