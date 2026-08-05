# Scenario

**Feature**: `--color` forces gray SGR on meta fields even with pipe writers

```
Harness -> list --color <root>
  -> gray (\x1b[90m) on leaf count / L2:L3 / summary; path plain
```

## Steps

1. Write one-leaf root.
2. Args = `list --color <root>`.

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
	req.Args = listArgs([]string{"--color"}, root)
	return nil
}
```
