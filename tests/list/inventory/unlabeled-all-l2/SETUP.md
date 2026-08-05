# Scenario

**Feature**: all unlabeled leaves → L2:L3=N:0 (100.0%/0.0%) and unlabeled=N

```
Harness -> 3 unlabeled leaves
  -> list <root>
  -> L2:L3=3:0 (100.0%/0.0%) unlabeled=3
```

## Steps

1. Write root with leaves a,b,c unlabeled.
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
	writeLabeledLeaves(t, root, []string{"a", "b", "c"})
	req.FixtureDir = root
	req.Roots = []string{root}
	req.Args = listArgs(nil, root)
	return nil
}
```
