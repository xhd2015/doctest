# Scenario

**Feature**: `doctest cache` on an empty/missing doctest root still prints home + root

```
empty CacheHome (no doctest/ yet)
  -> doctest cache
  -> Cache home + Doctest root; empty/0 buckets/0B
```

## Preconditions

- Parent info Setup set CacheHome to t.TempDir().
- DoctestRoot path is computed but directory is **not** created.

## Steps

1. Confirm DoctestRoot does not exist.
2. Set Args to `cache`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	ensureCacheHome(t, req)
	if pathExists(req.DoctestRoot) {
		t.Fatalf("precondition: DoctestRoot must not exist yet: %s", req.DoctestRoot)
	}
	req.Args = []string{"cache"}
	return nil
}
```
