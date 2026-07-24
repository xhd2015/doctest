# Scenario

**Feature**: two Go modules exist: mod-a has DOCTEST.md + SETUP.md, mod-b has its own SETUP.md

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- Two Go modules exist: mod-a has DOCTEST.md + SETUP.md, mod-b has its own SETUP.md.
- Running on a sub-dir of mod-b should not use mod-a's DOCTEST.md.

## Steps
1. Create mod-a (go.mod, tests/DOCTEST.md, tests/SETUP.md, leaf-a).
2. Create mod-b (go.mod, SETUP.md with types+Run, sub/leaf-b).
3. Run `doctest test <mod-b>/sub/leaf-b`.

```go
import (
    "os"
    "path/filepath"
    "testing"

    "github.com/xhd2015/doctest/libdoc/testtree"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    tmp := t.TempDir()
    bt := string(rune(96))
    d := bt + bt + bt

    modA := filepath.Join(tmp, "mod-a")
    modATests := filepath.Join(modA, "tests")
    os.MkdirAll(modATests, 0755)
    os.WriteFile(filepath.Join(modA, "go.mod"), []byte("module mod-a\n\ngo 1.21\n"), 0644)
    modADoctest := testtree.MinimalDOCTEST("import \"testing\"\n\ntype RequestA struct{}\ntype ResponseA struct{}\n\nfunc Run(t *testing.T, d *session.Doctest, req *RequestA) (*ResponseA, error) { return &ResponseA{}, nil }")
    os.WriteFile(filepath.Join(modATests, "DOCTEST.md"), []byte(modADoctest), 0644)
    leafA := filepath.Join(modATests, "leaf-a")
    os.MkdirAll(leafA, 0755)
    os.WriteFile(filepath.Join(leafA, "SETUP.md"), []byte(d+"go\nfunc Setup(t *testing.T, d *session.Doctest, req *RequestA) error { _ = req; return nil }\n"+d+"\n"), 0644)
    os.WriteFile(filepath.Join(leafA, "ASSERT.md"), []byte(d+"go\nfunc Assert(t *testing.T, d *session.Doctest, req *RequestA, resp *ResponseA, err error) {}\n"+d+"\n"), 0644)

    modB := filepath.Join(tmp, "mod-b")
    os.MkdirAll(modB, 0755)
    os.WriteFile(filepath.Join(modB, "go.mod"), []byte("module mod-b\n\ngo 1.21\n"), 0644)
    os.WriteFile(filepath.Join(modB, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(testtree.MinimalRunGo())), 0644)

    subDir := filepath.Join(modB, "sub")
    os.MkdirAll(subDir, 0755)
    leafSetup := d+"go\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { _ = req; return nil }\n"+d+"\n"
    os.WriteFile(filepath.Join(subDir, "SETUP.md"), []byte(leafSetup), 0644)

    leafB := filepath.Join(subDir, "leaf-b")
    os.MkdirAll(leafB, 0755)
    os.WriteFile(filepath.Join(leafB, "SETUP.md"), []byte(leafSetup), 0644)
    os.WriteFile(filepath.Join(leafB, "ASSERT.md"), []byte(d+"go\nfunc Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}\n"+d+"\n"), 0644)

    req.Args = []string{"test", leafB}
    return nil
}
```
