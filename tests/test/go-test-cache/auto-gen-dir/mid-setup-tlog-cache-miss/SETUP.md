# Scenario

**Feature**: intermediate Setup t.Log tag change is live → cache miss

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- mid-a Setup uses t.Log("…"); ASSERT empty is fine (Log is a side effect).
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
