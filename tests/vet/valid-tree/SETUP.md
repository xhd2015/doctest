# Scenario

**Feature**: a minimal valid doctest tree is available

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- A minimal valid doctest tree is available.

## Steps
1. Run `doctest vet <dir>`.

```go
import (
    "github.com/xhd2015/doctest/libdoc/testtree"
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    dir := t.TempDir()
    if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.VetDOCTEST()), 0644); err != nil {
        t.Fatal(err)
    }
    if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte("# Scenario\n\n**Feature**: minimal test setup\n\n\x60\x60\x60\n# minimal pipeline\nsystem -> run\n\x60\x60\x60\n\n## Setup\nsetup\n"), 0644); err != nil {
        t.Fatal(err)
    }
    req.Args = []string{"vet", dir}
    return nil
}
```
