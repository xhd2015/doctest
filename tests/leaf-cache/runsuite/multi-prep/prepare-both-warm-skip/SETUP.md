# Scenario

**Feature**: when both trees are warm, PreparePassPlan skip list has both tree-qualified identities

```
PutPass(keyA); PutPass(keyB)
PreparePassPlan([A/leaf, B/leaf], skipEnabled=true)
  -> len(Skip)==2
  -> Skip contains FormatLeafIdentity(A,"leaf") and FormatLeafIdentity(B,"leaf")
```

## Preconditions

- PrepWarm=`both`; SkipEnabled=true.
- Relative leaf path is `leaf` in both trees.

## Steps

1. Seed store for both trees.
2. Op=`multi_prep_prepare`.
3. Assert two skip identities, both tree-qualified and distinct.

## Context

- Core multi-prep contract: skip list never uses bare relpath alone.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "multi_prep_prepare"
	req.PrepWarm = "both"
	req.SkipEnabled = true
	req.SkipEnabledSet = true
	if err := seedWarmStore(t, req); err != nil {
		return err
	}
	return nil
}
```
