# Scenario

**Feature**: unread leaf SETUP WorkDir write may stay cached (DCE)

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- Leaf Setup sets WorkDir; ASSERT empty.
- Multi-run harness from parent.

## Steps
1. Edit leaf SETUP WorkDir tag; expect **still cached**.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg = multiRunCfg{
        TestDir:       createTempTestProjectLeafWorkDirDead(t, "mytest"),
        ModifyFile:    "simple/SETUP.md",
        ModifyContent: modifiedSetupContent("leaf-dead-v2"),
    }
    doMultiRun(t, req)
    return nil
}
```
