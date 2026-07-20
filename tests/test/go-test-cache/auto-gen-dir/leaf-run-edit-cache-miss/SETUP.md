# Scenario

**Feature**: DOCTEST Run body change uses t.Log (live side effect) → cache miss

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- Run calls `t.Log(tag)` so the string is not DCE'd from the binary.
- Multi-run harness from parent.

## Steps
1. Start with `t.Log("run-v1")`; edit to `t.Log("run-edited")` (meaningful → cache miss).

```go
import (
    "testing"

    "github.com/xhd2015/doctest/libdoc/testtree"
)

func Setup(t *testing.T, req *Request) error {
    testDir := createTempTestProjectOpts(t, "mytest", treeOpts{
        RunCode: runCodeWithLog("run-v1"),
    })
    cfg = multiRunCfg{
        TestDir:       testDir,
        ModifyFile:    "DOCTEST.md",
        ModifyContent: testtree.MinimalDOCTEST(doctestBody(modifiedRunCode())),
    }
    doMultiRun(t, req)
    return nil
}
```
