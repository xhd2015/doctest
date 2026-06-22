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
- Build a group directory with multiple leaves.

## Steps
1. Create the same test tree.
2. Run `doctest build <treeRoot>/group-a`.

```go
import (
    "os"
    "path/filepath"
    "testing"

    "github.com/xhd2015/doctest/libdoc/testtree"
)

func Setup(t *testing.T, req *Request) error {
    treeRoot := t.TempDir()
    bt := string(rune(96))
    d := bt + bt + bt

    testtree.WriteFile(t, treeRoot, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))

    groupSetup := d+"go\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"+d+"\n"
    leafSetup := d+"go\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"+d+"\n"
    leafAssert := d+"go\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"+d+"\n"

    ga := filepath.Join(treeRoot, "group-a")
    os.MkdirAll(ga, 0755)
    os.WriteFile(filepath.Join(ga, "SETUP.md"), []byte(groupSetup), 0644)

    leaf1 := filepath.Join(ga, "leaf-1")
    os.MkdirAll(leaf1, 0755)
    os.WriteFile(filepath.Join(leaf1, "SETUP.md"), []byte(leafSetup), 0644)
    os.WriteFile(filepath.Join(leaf1, "ASSERT.md"), []byte(leafAssert), 0644)

    leaf2 := filepath.Join(ga, "leaf-2")
    os.MkdirAll(leaf2, 0755)
    os.WriteFile(filepath.Join(leaf2, "SETUP.md"), []byte(leafSetup), 0644)
    os.WriteFile(filepath.Join(leaf2, "ASSERT.md"), []byte(leafAssert), 0644)

    req.Args = []string{"build", ga}
    return nil
}
```
