# Scenario

**Feature**: the `validate` command has been renamed to `vet`

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- The `validate` command has been renamed to `vet`.

## Steps
1. Run `doctest validate <dir>`.

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
    req.Args = []string{"validate", dir}
    return nil
}
```
