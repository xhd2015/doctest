# Scenario

**Feature**: override roots outside the main doctest tree appear in clean plan

```
# outside DOCTEST_LEAF_CACHE (injected on Request, not process env)
main: CacheHome/doctest/...
outside: LeafCache abs path (not under doctest root)
  -> cache --clean --dry-run
  -> would remove both paths
```

## Preconditions

- Override paths are absolute and **outside** `DoctestRoot`.
- Injected via `Request.LeafCache` / `Request.MetricsRoot` — never process Setenv.

## Steps

1. Grouping Setup ensures CacheHome.
2. Leaf creates outside root and sets clean dry-run Args.

## Context

- Grouping only; one recommended override leaf for MVP.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	ensureCacheHome(t, req)
	return nil
}
```
