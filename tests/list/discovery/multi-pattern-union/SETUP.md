# Scenario

**Feature**: multiple patterns form a sorted, deduped union

```
Harness -> roots a, b
  -> list a b a
  -> two body lines sorted; a not duplicated
```

## Preconditions

- Two distinct roots under temp base.

## Steps

1. Write roots `a` and `b`.
2. Args = `list <a> <b> <a>` (duplicate a).

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	base := t.TempDir()
	a := filepath.Join(base, "a")
	b := filepath.Join(base, "b")
	writeLabeledLeaves(t, a, []string{"leaf"})
	writeLabeledLeaves(t, b, []string{"leaf"})
	req.FixtureDir = base
	req.Roots = []string{a, b}
	req.Args = listArgs(nil, a, b, a)
	return nil
}
```
