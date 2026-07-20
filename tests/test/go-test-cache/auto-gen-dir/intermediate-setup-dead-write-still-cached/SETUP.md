# Scenario

**Feature**: unread intermediate SETUP WorkDir write may leave go testcache hit

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- L1 mid Setup sets WorkDir; leaf ASSERT is empty (never reads WorkDir).
- Multi-run harness from parent.
- Go may DCE the write → binary content ID stable → still cached.

## Steps
1. Edit `mid-a/SETUP.md` WorkDir only; expect **still cached** (unused change).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    // Dead tree: intermediate writes WorkDir; ASSERT does not observe it.
    cfg = multiRunCfg{
        TestDir:       createL1ProjectDead(t, "mytest"),
        ModifyFile:    "mid-a/SETUP.md",
        ModifyContent: modifiedSetupContent("dead-l1-v2"),
    }
    doMultiRun(t, req)
    return nil
}
```
