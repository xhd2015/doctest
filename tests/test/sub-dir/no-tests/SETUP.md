# Scenario

**Feature**: a doc-style test tree exists

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- A doc-style test tree exists.
- A sub-directory has no ASSERT.md files (no test cases).

## Steps
1. Create a test tree with a no-leaf-dir that has a SETUP.md but no children with ASSERT.md.
2. Run `doctest test <treeRoot>/no-leaf-dir`.

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

    noLeafDir := filepath.Join(treeRoot, "no-leaf-dir")
    os.MkdirAll(noLeafDir, 0755)
    os.WriteFile(filepath.Join(noLeafDir, "SETUP.md"), []byte(
        bt+"go\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }\n"+bt+"\n"), 0644)

    req.Args = []string{"test", noLeafDir}
    return nil
}
```
