# Scenario

**Feature**: `--no-color` forces plain output (no ANSI)

```
Harness -> list --no-color <root>
  -> no ESC in stdout
```

## Steps

1. Write one-leaf root.
2. Args = `list --no-color <root>`.

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
	req.Args = listArgs([]string{"--no-color"}, root)
	return nil
}
```
