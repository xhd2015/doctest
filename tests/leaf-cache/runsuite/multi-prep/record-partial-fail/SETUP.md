# Scenario

**Feature**: RecordPasses stores only non-failed identities on partial suite fail

```
keys = {idA: keyA, idB: keyB}
failed = {idA: true}
RecordPasses(store, keys, failed, allPassed=false)
  -> GetPass(keyA)=false
  -> GetPass(keyB)=true
```

## Preconditions

- Twin trees; empty store before record (no pre-warm required).
- FailedSide=`a` marks treeA's identity as failed.
- Op=`multi_prep_record`.

## Steps

1. Build keys for both leaves via ComputeLeafKey.
2. RecordPasses with failed={identityA}.
3. Assert GetPass only for the non-failed leaf.

## Context

- Locked product rule: PutPass on pass even if suite partially fails; fail never
  PutPass. Multi-prep uses tree-qualified identity keys in the failed map.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	t.Helper()
	req.Op = "multi_prep_record"
	req.FailedSide = "a"
	req.AllPassed = false
	return nil
}
```
