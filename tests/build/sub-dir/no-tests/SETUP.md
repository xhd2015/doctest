# Scenario

**Feature**: a doc-style test tree exists

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Preconditions
- A doc-style test tree exists.
- A sub-directory has no ASSERT.md files.

## Steps
1. Create a test tree with a no-leaf-dir that has a SETUP.md but no ASSERT.md descendants.
2. Run `doctest build <treeRoot>/no-leaf-dir`.

```go
import (
    "os"
    "path/filepath"
    "testing"

    "github.com/xhd2015/doctest/libdoc/testtree"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    treeRoot := t.TempDir()
    bt := string(rune(96))
    d := bt + bt + bt

    testtree.WriteFile(t, treeRoot, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))

    noLeafDir := filepath.Join(treeRoot, "no-leaf-dir")
    os.MkdirAll(noLeafDir, 0755)
    os.WriteFile(filepath.Join(noLeafDir, "SETUP.md"), []byte(
        d+"go\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }\n"+d+"\n"), 0644)

    req.Args = []string{"build", noLeafDir}
    return nil
}
```
