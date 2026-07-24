# Scenario

**Feature**: a doc-style test tree exists with multiple groups

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A doc-style test tree exists with multiple groups.
- Run on a group directory.

## Steps
1. Create the same test tree as run-specific-leaf.
2. Run `doctest test <treeRoot>/group-a`.

```go
import (
    "os"
    "path/filepath"
    "testing"

    "github.com/xhd2015/doctest/libdoc/testtree"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    treeRoot := t.TempDir()
    bt := "\x60\x60\x60"

    testtree.WriteFile(t, treeRoot, "DOCTEST.md", testtree.MinimalDOCTEST(testtree.MinimalRunGo()))

    groupSetup := bt+"go\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }\n"+bt+"\n"
    leafSetup := bt+"go\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }\n"+bt+"\n"
    leafAssert := bt+"go\nfunc Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}\n"+bt+"\n"

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

    gb := filepath.Join(treeRoot, "group-b")
    os.MkdirAll(gb, 0755)
    os.WriteFile(filepath.Join(gb, "SETUP.md"), []byte(groupSetup), 0644)

    leaf3 := filepath.Join(gb, "leaf-3")
    os.MkdirAll(leaf3, 0755)
    os.WriteFile(filepath.Join(leaf3, "SETUP.md"), []byte(leafSetup), 0644)
    os.WriteFile(filepath.Join(leaf3, "ASSERT.md"), []byte(leafAssert), 0644)

    req.Args = []string{"test", ga}
    return nil
}
```
