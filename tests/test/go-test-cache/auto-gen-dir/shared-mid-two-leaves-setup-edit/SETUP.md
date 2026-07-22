# Scenario

**Feature**: shared intermediate SETUP edit with observed WorkDir

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- mid-a shared by two leaves; ASSERT observes WorkDir.
- Multi-run harness from parent.

## Steps
1. Edit shared `mid-a/SETUP.md` (on both leaves' spines → both keys miss → 0 Cached).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    cfg := multiRunCfg{
        TestDir:       createSharedMidTwoLeavesProject(t, "mytest"),
        ModifyFile:    "mid-a/SETUP.md",
        ModifyContent: modifiedSetupContent("modified-shared-mid"),
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
