# Scenario

**Feature**: `base/...` discovers multiple sibling roots, sorted

```
Harness -> base/{alpha,beta} each with DOCTEST.md
  -> list base/...
  -> two body lines (sorted) + summary
```

## Preconditions

- Temp base with two sibling roots `alpha` (1 leaf) and `beta` (1 leaf).

## Steps

1. Write sibling roots under base.
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
	alpha := filepath.Join(base, "alpha")
	beta := filepath.Join(base, "beta")
	writeLabeledLeaves(t, alpha, []string{"leaf"})
	writeLabeledLeaves(t, beta, []string{"leaf"})
	req.FixtureDir = base
	req.Roots = []string{alpha, beta}
	// path/... form without relying on process cwd
	req.Args = listArgs(nil, filepath.ToSlash(base)+"/...")
	return nil
}
```
