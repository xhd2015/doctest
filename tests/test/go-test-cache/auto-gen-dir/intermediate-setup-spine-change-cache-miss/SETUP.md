# Scenario

**Feature**: intermediate SETUP Go change busts leaf-cache (spine content hash)

```
# mid-a SETUP is on the leaf spine
# even unread WorkDir writes change the spine hash -> miss
doctest test <fixture> after mid SETUP edit -> 0 Cached
```

## Preconditions
- L1 mid Setup sets WorkDir; leaf ASSERT is empty (never reads WorkDir).
- Leaf-cache keys spine Go text (not go DCE).
- Multi-run harness from parent.

## Steps
1. Edit `mid-a/SETUP.md` WorkDir only; expect **0 Cached**.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg := multiRunCfg{
        TestDir:       createL1ProjectDead(t, "mytest"),
        ModifyFile:    "mid-a/SETUP.md",
        ModifyContent: modifiedSetupContent("dead-l1-v2"),
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
