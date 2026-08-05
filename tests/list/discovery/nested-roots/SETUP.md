# Scenario

**Feature**: nested DOCTEST roots listed separately; parent leaf count excludes nested tree

```
Harness -> parent (leaf + nested/child root with its own leaf)
  -> list base/...
  -> parent leaves exclude nested ASSERT; child listed separately
```

## Preconditions

- Parent root owns one leaf; nested child root owns one leaf.

## Steps

1. Write parent with leaf `own` and nested `nested/` DOCTEST root with leaf `child-leaf`.
2. Args = `list <base>/...`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	base := t.TempDir()
	parent := filepath.Join(base, "parent")
	writeLabeledLeaves(t, parent, []string{"own"})
	child := filepath.Join(parent, "nested")
	writeLabeledLeaves(t, child, []string{"child-leaf"})
	req.FixtureDir = base
	req.Roots = []string{parent, child}
	req.Args = listArgs(nil, filepath.ToSlash(base)+"/...")
	return nil
}
```
