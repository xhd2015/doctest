# Scenario

**Feature**: intermediate `_ = "tag"` string change busts leaf-cache (spine hash)

```
# discarded string constants still appear in SETUP Go text
# leaf-cache hashes spine content -> tag v1→v2 misses
doctest test <fixture> after discard-tag edit -> 0 Cached
```

## Preconditions
- mid-a Setup only discards a string constant; ASSERT empty.
- Spine content hashing (not go DCE) decides leaf-cache hits.
- Multi-run harness from parent.

## Steps
1. Edit mid-a SETUP discard tag v1→v2; expect **0 Cached**.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg = multiRunCfg{
        TestDir:       createL1ProjectDeadDiscard(t, "mytest"),
        ModifyFile:    "mid-a/SETUP.md",
        ModifyContent: modifiedDiscardStringSetup("discard-v2"),
    }
    doMultiRun(t, req)
    return nil
}
```
