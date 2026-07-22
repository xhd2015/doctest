# Scenario

**Feature**: any root SETUP Go change busts leaf-cache (spine content hash)

```
# root SETUP Go is on every leaf's spine
# even "unread" WorkDir writes change the spine hash -> miss
doctest test <fixture> after SETUP.md edit -> 0 Cached
```

## Preconditions
- Root SETUP sets WorkDir; leaf ASSERT is empty (value unused at runtime).
- Leaf-cache keys the **spine Go text**, not go binary DCE / testlog liveness.
- Multi-run harness from parent.

## Steps
1. Edit root `SETUP.md` WorkDir tag only; expect **0 Cached** (spine change).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg := multiRunCfg{
        TestDir:       createTempTestProjectRootWorkDirDead(t, "mytest", "root-dead-v1"),
        ModifyFile:    "SETUP.md",
        ModifyContent: modifiedSetupContent("root-dead-v2"),
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
