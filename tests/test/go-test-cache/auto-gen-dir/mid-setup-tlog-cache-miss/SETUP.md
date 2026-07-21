# Scenario

**Feature**: intermediate Setup t.Log tag change busts leaf-cache (spine hash)

```
# mid SETUP Go is on the leaf spine; tag string is part of spine text
doctest test <fixture> after t.Log tag edit -> 0 Cached
```

## Preconditions
- mid-a Setup uses t.Log("…"); spine content hashing includes the tag.
- Multi-run harness from parent.

## Steps
1. Edit mid-a SETUP t.Log tag v1→v2; expect **0 Cached**.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg = multiRunCfg{
        TestDir:       createL1ProjectMidTLog(t, "mytest"),
        ModifyFile:    "mid-a/SETUP.md",
        ModifyContent: modifiedTLogSetup("mid-tlog-v2"),
    }
    doMultiRun(t, req)
    return nil
}
```
