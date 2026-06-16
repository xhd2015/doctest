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
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    dir := t.TempDir()
    if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte("# tests\n"), 0644); err != nil {
        t.Fatal(err)
    }
    if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte("setup\n"), 0644); err != nil {
        t.Fatal(err)
    }
    req.Args = []string{"vet", dir}
    return nil
}
```
