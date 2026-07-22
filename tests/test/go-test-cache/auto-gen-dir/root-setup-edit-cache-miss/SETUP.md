# Scenario

**Feature**: root SETUP WorkDir is observed by leaf ASSERT

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- Root SETUP sets WorkDir; leaf ASSERT requires non-empty WorkDir.
- Multi-run harness from parent.

## Steps
1. Edit root `SETUP.md` WorkDir between runs (meaningful → cache miss).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg := multiRunCfg{
        TestDir:       createTempTestProjectRootWorkDir(t, "mytest", "root-v1"),
        ModifyFile:    "SETUP.md",
        ModifyContent: modifiedSetupContent("modified-root-setup"),
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
