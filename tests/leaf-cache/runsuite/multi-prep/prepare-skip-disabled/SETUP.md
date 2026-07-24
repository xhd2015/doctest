# Scenario

**Feature**: when skip is disabled, PreparePassPlan returns empty Skip even if store is warm

```
PutPass(keyA); PutPass(keyB)
PreparePassPlan(..., skipEnabled=false)
  -> Skip empty
  -> Keys still populated (for later record)
```

## Preconditions

- PrepWarm=`both`; SkipEnabled=false (models `-count` / `-a`).

## Steps

1. Seed both trees warm.
2. Op=`multi_prep_prepare` with SkipEnabled=false.
3. Assert SkipPaths empty; keys/identities still present.

## Context

- Mirrors `leafcache.SkipEnabled` product rule at the multi-prep seam.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "multi_prep_prepare"
	req.PrepWarm = "both"
	req.SkipEnabled = false
	req.SkipEnabledSet = true
	if err := seedWarmStore(t, req); err != nil {
		return err
	}
	return nil
}
```
