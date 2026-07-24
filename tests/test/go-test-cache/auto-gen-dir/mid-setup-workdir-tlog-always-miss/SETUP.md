# Scenario

**Feature**: mid WorkDir write + ASSERT always t.Log(WorkDir) → cache miss

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- mid-a Setup sets WorkDir; leaf ASSERT always t.Log(req.WorkDir) (strong live).
- Multi-run harness from parent.

## Steps
1. Edit mid-a SETUP WorkDir tag; expect **0 Cached**.

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    cfg := multiRunCfg{
        TestDir:       createL1ProjectWorkDirTLogAssert(t, "mytest"),
        ModifyFile:    "mid-a/SETUP.md",
        ModifyContent: modifiedSetupContent("mid-workdir-tlog-v2"),
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
