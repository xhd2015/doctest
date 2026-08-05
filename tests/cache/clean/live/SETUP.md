# Scenario

**Feature**: `--clean` removes the entire doctest root tree

```
seed buckets under DoctestRoot
  -> doctest cache --clean
  -> Removed <abs> (<size>); DoctestRoot gone
```

## Preconditions

- Injectable CacheHome; seeded content so remove has something to wipe.

## Steps

1. Seed two buckets so the tree is non-empty.
2. Set Args to `cache --clean`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	ensureCacheHome(t, req)
	seedBucket(t, req, "leaf-cache", seedBytes(256))
	seedBucket(t, req, "mapping-gen", seedBytes(256))
	if !pathExists(req.DoctestRoot) {
		t.Fatalf("precondition: DoctestRoot should exist after seed: %s", req.DoctestRoot)
	}
	req.Args = []string{"cache", "--clean"}
	return nil
}
```
