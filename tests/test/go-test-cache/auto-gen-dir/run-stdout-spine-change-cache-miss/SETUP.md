# Scenario

**Feature**: DOCTEST Run Go change busts leaf-cache (spine content hash)

```
# Run is part of root DOCTEST Go on the spine
# Stdout string swap changes spine hash even if ASSERT only checks non-empty
doctest test <fixture> after DOCTEST Run edit -> 0 Cached
```

## Preconditions
- Run returns Stdout "run-v1"; ASSERT only checks Stdout != "".
- Leaf-cache keys spine Go text (not go DCE of string constants).
- Multi-run harness from parent.

## Steps
1. Edit DOCTEST Run Stdout to "run-edited"; expect **0 Cached**.

```go
import (
    "os"
    "path/filepath"
    "testing"

    "github.com/xhd2015/doctest/libdoc/testtree"
)

func Setup(t *testing.T, req *Request) error {
    testDir := createTempTestProjectOpts(t, "mytest", treeOpts{
        RunCode: runCodeWithStdout("run-v1"),
    })
    if err := os.WriteFile(filepath.Join(testDir, "simple", "ASSERT.md"), []byte(leafAssertStdoutNonEmpty()), 0644); err != nil {
        t.Fatalf("write assert: %v", err)
    }
    cfg := multiRunCfg{
        TestDir:    testDir,
        ModifyFile: "DOCTEST.md",
        ModifyContent: testtree.MinimalDOCTEST(doctestBody(runCodeWithStdout("run-edited"))),
    }
    doMultiRun(t, req, cfg)
    return nil
}
```
