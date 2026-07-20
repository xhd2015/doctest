# Scenario

**Feature**: L1 intermediate SETUP WorkDir is observed by leaf ASSERT

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- L1 mid-a Setup sets WorkDir; leaf ASSERT requires non-empty WorkDir.
- Multi-run harness from parent.

## Steps
1. Edit only `mid-a/SETUP.md` WorkDir tag (meaningful → cache miss).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg = multiRunCfg{
        TestDir:       createL1Project(t, "mytest"),
        ModifyFile:    "mid-a/SETUP.md",
        ModifyContent: modifiedSetupContent("modified-l1-setup"),
    }
    doMultiRun(t, req)
    return nil
}
```
