# Scenario

**Feature**: a first run has completed successfully

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- Flat tree with leaf Setup setting WorkDir; ASSERT observes non-empty WorkDir (live code).
- Multi-run harness from parent.

## Steps
1. Edit leaf `simple/SETUP.md` WorkDir tag between runs (meaningful → cache miss).

```go
import (
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    cfg := multiRunCfg{
        TestDir:       createTempTestProjectObserveWorkDir(t, "mytest"),
        ModifyFile:    "simple/SETUP.md",
        ModifyContent: modifiedSetupContent("modified-leaf-setup"),
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
