# Scenario

**Feature**: L2-only intermediate SETUP edit with observed WorkDir

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- mid-a + mid-a/mid-b; leaf ASSERT observes WorkDir.
- Multi-run harness from parent.

## Steps
1. Edit only `mid-a/mid-b/SETUP.md` (spine ancestor → leaf-cache miss).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg = multiRunCfg{
        TestDir:       createL2Project(t, "mytest"),
        ModifyFile:    "mid-a/mid-b/SETUP.md",
        ModifyContent: modifiedSetupContent("modified-l2-setup"),
    }
    doMultiRun(t, req)
    return nil
}
```
