# Scenario

**Feature**: `doctest cache --clean` removes or previews the doctest cache tree

```
# clean matrix
seed CacheHome/doctest/...
  -> --dry-run alone        -> error (requires --clean)
  -> --clean --dry-run      -> [dry-run] would remove; no delete
  -> --clean                -> Removed; tree gone
```

## Preconditions

- Leaves isolate CacheHome and seed at least one bucket unless testing
  flag-only dry-run-without-clean.
- Parallel-safe: each leaf owns its temp tree.

## Steps

1. Grouping Setup ensures a temp CacheHome.
2. Leaves seed (when needed) and set clean-related Args.

## Context

- Grouping: shared CacheHome isolation default.

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
