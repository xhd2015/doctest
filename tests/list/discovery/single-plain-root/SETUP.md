# Scenario

**Feature**: list a single plain root path → one body line + summary

```
Harness -> write one root with 2 unlabeled leaves
  -> list <absRoot>
  -> one body line + summary footer
```

## Preconditions

- Temp root with two unlabeled ASSERT leaves.

## Steps

1. Create fixture root with leaves `a` and `b`.
2. Args = `list <absRoot>`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	root := filepath.Join(t.TempDir(), "tree")
	writeLabeledLeaves(t, root, []string{"a", "b"})
	req.FixtureDir = root
	req.Roots = []string{root}
	req.Args = listArgs(nil, root)
	return nil
}
```
