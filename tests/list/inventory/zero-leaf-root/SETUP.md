# Scenario

**Feature**: zero-leaf root (DOCTEST only) → leaves=0, L2:L3=0:0, no percent group

```
Harness -> root with DOCTEST.md only (no ASSERT)
  -> list <root>
  -> leaves=0 L2:L3=0:0 without (p2%/p3%)
```

## Steps

1. Write root DOCTEST only.
2. Args = `list <root>`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	root := filepath.Join(t.TempDir(), "empty-tree")
	writeRootDOCTEST(t, root)
	req.FixtureDir = root
	req.Roots = []string{root}
	req.Args = listArgs(nil, root)
	return nil
}
```
