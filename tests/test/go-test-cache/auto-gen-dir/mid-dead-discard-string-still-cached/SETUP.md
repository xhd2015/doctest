# Scenario

**Feature**: intermediate Setup `_ = "tag"` only may stay cached (DCE)

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- mid-a Setup only discards a string constant; ASSERT empty.
- Multi-run harness from parent.

## Steps
1. Edit mid-a SETUP discard tag v1→v2; expect **still cached**.

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
