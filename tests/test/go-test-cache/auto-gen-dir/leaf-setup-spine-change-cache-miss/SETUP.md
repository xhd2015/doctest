# Scenario

**Feature**: leaf SETUP Go change busts leaf-cache (spine content hash)

```
# leaf SETUP is on the spine
# even unread WorkDir writes change the spine hash -> miss
doctest test <fixture> after leaf SETUP edit -> 0 Cached
```

## Preconditions
- Leaf Setup sets WorkDir; ASSERT empty.
- Leaf-cache keys spine Go text (not go DCE).
- Multi-run harness from parent.

## Steps
1. Edit leaf SETUP WorkDir tag; expect **0 Cached**.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg := multiRunCfg{
        TestDir:       createTempTestProjectLeafWorkDirDead(t, "mytest"),
        ModifyFile:    "simple/SETUP.md",
        ModifyContent: modifiedSetupContent("leaf-dead-v2"),
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
