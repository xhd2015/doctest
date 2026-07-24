# Scenario

**Feature**: only the warm tree's identity appears in the multi-prep skip list

```
PutPass(keyA) only
PreparePassPlan([A/leaf, B/leaf], skipEnabled=true)
  -> Skip == [identity(A)]
  -> identity(B) not in Skip
```

## Preconditions

- PrepWarm=`a` (treeA only); SkipEnabled=true.

## Steps

1. Seed store for treeA only.
2. Op=`multi_prep_prepare`.
3. Assert single skip entry equals treeA identity.

## Context

- Prevents cold treeB from inheriting a skip because of shared relative path.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "multi_prep_prepare"
	req.PrepWarm = "a"
	req.SkipEnabled = true
	req.SkipEnabledSet = true
	if err := seedWarmStore(t, req); err != nil {
		return err
	}
	return nil
}
```
