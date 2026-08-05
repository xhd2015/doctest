# Scenario

**Feature**: `--clean --dry-run` previews removal and does not delete

```
seed buckets under DoctestRoot
  -> doctest cache --clean --dry-run
  -> [dry-run] would remove: <abs> (<size>); tree still exists
```

## Preconditions

- Injectable CacheHome; at least one seeded bucket so size is non-trivial.

## Steps

1. Seed `leaf-cache` with a small payload.
2. Set Args to `cache --clean --dry-run`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	ensureCacheHome(t, req)
	seedBucket(t, req, "leaf-cache", seedBytes(512))
	if !pathExists(req.DoctestRoot) {
		t.Fatalf("precondition: DoctestRoot should exist after seed: %s", req.DoctestRoot)
	}
	req.Args = []string{"cache", "--clean", "--dry-run"}
	return nil
}
```
