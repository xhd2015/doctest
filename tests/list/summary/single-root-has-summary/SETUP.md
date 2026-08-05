# Scenario

**Feature**: single root still emits summary footer and trailing newline

```
Harness -> one root one leaf
  -> list <root>
  -> body + blank + --- + totals + labels + trailing \n
```

## Steps

1. Write one-leaf root.
2. Args = `list <root>`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	root := filepath.Join(t.TempDir(), "one")
	writeLabeledLeaves(t, root, []string{"only"})
	req.FixtureDir = root
	req.Roots = []string{root}
	req.Args = listArgs(nil, root)
	return nil
}
```
