# Scenario

**Feature**: an invalid doctest tree has an assertion without local setup

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- An invalid doctest tree has an assertion without local setup.

## Steps
1. Run `doctest vet <dir>`.

```go
import (
    "github.com/xhd2015/doctest/libdoc/testtree"
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    dir := t.TempDir()
    if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.VetDOCTEST()), 0644); err != nil {
        t.Fatal(err)
    }
    if err := os.MkdirAll(filepath.Join(dir, "leaf"), 0755); err != nil {
        t.Fatal(err)
    }
    if err := os.WriteFile(filepath.Join(dir, "leaf", "ASSERT.md"), []byte("assert\n"), 0644); err != nil {
        t.Fatal(err)
    }
    req.Args = []string{"vet", dir}
    return nil
}
```
