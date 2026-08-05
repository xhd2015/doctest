# Scenario

**Feature**: mix of e2e and unlabeled → L3 count and percents

```
Harness -> 1 e2e leaf + 1 unlabeled
  -> list <root>
  -> L2:L3=1:1 (50.0%/50.0%); e2e=1 unlabeled=1
```

## Steps

1. Write leaves `fast` (unlabeled) and `e2e-leaf|e2e`.
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
	writeLabeledLeaves(t, root, []string{"fast", "e2e-leaf|e2e"})
	req.FixtureDir = root
	req.Roots = []string{root}
	req.Args = listArgs(nil, root)
	return nil
}
```
