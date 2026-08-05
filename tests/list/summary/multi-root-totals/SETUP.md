# Scenario

**Feature**: multi-root summary totals and labels equal the sum of body rows

```
Harness -> root-a: 2 unlabeled; root-b: 1 e2e
  -> list a b
  -> roots=2 leaves=3 L2:L3=2:1; labels sum e2e + unlabeled
```

## Steps

1. Write two roots with different inventory.
2. Args = `list <a> <b>`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	base := t.TempDir()
	a := filepath.Join(base, "root-a")
	b := filepath.Join(base, "root-b")
	writeLabeledLeaves(t, a, []string{"x", "y"})
	writeLabeledLeaves(t, b, []string{"z|e2e"})
	req.FixtureDir = base
	req.Roots = []string{a, b}
	req.Args = listArgs(nil, a, b)
	return nil
}
```
