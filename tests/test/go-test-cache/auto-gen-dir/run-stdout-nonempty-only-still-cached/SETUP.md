# Scenario

**Feature**: Run Stdout string swap with only non-empty ASSERT may stay cached (DCE)

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code
```

## Preconditions
- Run returns Stdout "run-v1"; ASSERT only checks Stdout != "".
- Go can DCE specific string constants when only emptiness is observed.
- Multi-run harness from parent.

## Steps
1. Edit DOCTEST Run Stdout to "run-edited"; expect **still cached**.

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
    cfg = multiRunCfg{
        TestDir:    testDir,
        ModifyFile: "DOCTEST.md",
        ModifyContent: testtree.MinimalDOCTEST(doctestBody(runCodeWithStdout("run-edited"))),
    }
    doMultiRun(t, req)
    return nil
}
```
