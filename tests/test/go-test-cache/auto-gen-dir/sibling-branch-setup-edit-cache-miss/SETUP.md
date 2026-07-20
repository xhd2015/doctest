# Scenario

**Feature**: sibling-branch intermediate SETUP edit with observed WorkDir

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- mid-a/leaf-x and other/leaf-z; ASSERT observes WorkDir on both leaves.
- Multi-run harness from parent.

## Steps
1. Edit only `other/SETUP.md` (meaningful → suite cache miss under unified package).

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    testDir := createTempCustomProjectOpts(t, "mytest", treeOpts{ObserveWorkDir: true},
        []string{"mid-a", "other"},
        []string{"mid-a/leaf-x", "other/leaf-z"},
    )
    cfg = multiRunCfg{
        TestDir:       testDir,
        ModifyFile:    "other/SETUP.md",
        ModifyContent: modifiedSetupContent("modified-sibling-branch"),
    }
    doMultiRun(t, req)
    return nil
}
```
