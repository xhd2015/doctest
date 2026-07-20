# Scenario

**Feature**: unread root SETUP WorkDir write may leave go testcache hit

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- Root SETUP sets WorkDir; leaf ASSERT is empty.
- Multi-run harness from parent.

## Steps
1. Edit root `SETUP.md` WorkDir only; expect **still cached**.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg = multiRunCfg{
        TestDir:       createTempTestProjectRootWorkDirDead(t, "mytest", "root-dead-v1"),
        ModifyFile:    "SETUP.md",
        ModifyContent: modifiedSetupContent("root-dead-v2"),
    }
    doMultiRun(t, req)
    return nil
}
```
