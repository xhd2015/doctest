# Scenario

**Feature**: L3-only intermediate SETUP edit with observed WorkDir

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- Deep chain; leaf ASSERT observes WorkDir.
- Multi-run harness from parent.

## Steps
1. Edit only `mid-a/mid-b/mid-c/SETUP.md` (spine ancestor → leaf-cache miss).
```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    cfg := multiRunCfg{
        TestDir:       createDeepChainProject(t, "mytest"),
        ModifyFile:    "mid-a/mid-b/mid-c/SETUP.md",
        ModifyContent: modifiedSetupContent("modified-l3-setup"),
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
