# Scenario

**Feature**: DOCTEST Run body change busts leaf-cache (spine hash)

```
# root DOCTEST Run Go is on the leaf spine
doctest test <fixture> after Run edit -> 0 Cached
```

## Preconditions
- Run calls `t.Log(tag)`; tag change is a spine Go edit.
- Multi-run harness from parent.

## Steps
1. Start with `t.Log("run-v1")`; edit to `t.Log("run-edited")` (spine miss).

```go
import (
    "testing"

    "github.com/xhd2015/doctest/libdoc/testtree"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    testDir := createTempTestProjectOpts(t, "mytest", treeOpts{
        RunCode: runCodeWithLog("run-v1"),
    })
    cfg := multiRunCfg{
        TestDir:       testDir,
        ModifyFile:    "DOCTEST.md",
        ModifyContent: testtree.MinimalDOCTEST(doctestBody(modifiedRunCode())),
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
