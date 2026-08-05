# Scenario

**Feature**: `doctest cache` lists seeded bucket names with human sizes and Total

```
seed leaf-cache (~1.2K) + mapping-gen (~4K)
  -> doctest cache
  -> bucket names on stdout; Total non-zero; human units
```

## Preconditions

- Injectable CacheHome from parent info Setup.
- At least two first-level buckets with known payload sizes.

## Steps

1. Seed `leaf-cache` with 1200 bytes and `mapping-gen` with 4096 bytes.
2. Set Args to `cache`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	ensureCacheHome(t, req)
	// ~1.2K and 4.0K human sizes (exact formatting is product choice).
	seedBucket(t, req, "leaf-cache", seedBytes(1200))
	seedBucket(t, req, "mapping-gen", seedBytes(4096))
	req.Args = []string{"cache"}
	return nil
}
```
