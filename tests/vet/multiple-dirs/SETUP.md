# Scenario

**Feature**: two independent valid doctest trees exist in different directories

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- Two independent valid doctest trees exist in different directories.

## Steps
1. Create two separate temp directories each with a valid doctest tree.
2. Run `doctest vet <dir1> <dir2>`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	var dirs []string
	for i := 0; i < 2; i++ {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte("# tests\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte("## Setup\n"), 0644); err != nil {
			t.Fatal(err)
		}
		dirs = append(dirs, dir)
	}
	req.Args = append([]string{"vet"}, dirs...)
	return nil
}
```
