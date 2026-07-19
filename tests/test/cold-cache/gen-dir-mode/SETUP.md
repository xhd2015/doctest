# Scenario

**Feature**: `--cold-cache` chooses and protects the mapping-gen root

```
# gen-dir resolution under cold-cache
omit --gen-dir -> $CacheHome/doctest/mapping-gen-cold (wipe on startup only)
--gen-dir X outside warm -> wipe X on startup, use X, leftover after finish
--gen-dir equal/under warm mapping-gen -> error, no wipe
```

## Preconditions

- Parent `cold-cache/SETUP.md` provides fixture and cache-sandbox helpers.
- Split factor at this level: **gen-dir resolution mode** (auto / allowed explicit / rejected).

## Steps

1. Descendant leaves set `--cold-cache` and optional `--gen-dir`, then assert wipe/reject policy.

```go
import (
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	// Gen-dir policy leaves all invoke --cold-cache against a tiny fixture that
	// must generate + go test; keep a generous timeout even if an ancestor lowered it.
	if req.Timeout < 120*time.Second {
		req.Timeout = 120 * time.Second
	}
	return nil
}
```
