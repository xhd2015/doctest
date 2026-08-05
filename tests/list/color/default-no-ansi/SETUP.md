# Scenario

**Feature**: default capture (pipe writers) has no ANSI

```
Harness -> list <root>  # no color flags; buffer writers
  -> stdout without ESC sequences
```

## Steps

1. Write one-leaf root.
2. Args = `list <root>` (no color flags).

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	root := filepath.Join(t.TempDir(), "tree")
	writeLabeledLeaves(t, root, []string{"a"})
	req.FixtureDir = root
	req.Roots = []string{root}
	req.Args = listArgs(nil, root)
	return nil
}
```
