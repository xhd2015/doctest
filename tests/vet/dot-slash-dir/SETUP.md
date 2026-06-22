# Scenario

**Feature**: a valid doctest tree exists in the current working directory

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- A valid doctest tree exists in the current working directory.

## Steps
1. Create a minimal valid doctest tree in a temp directory.
2. Run `doctest vet ./` with WorkDir set to that directory.

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
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte("# Scenario\n\n**Feature**: minimal test setup\n\n\x60\x60\x60\n# minimal pipeline\nsystem -> run\n\x60\x60\x60\n\n## Setup\n"), 0644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = dir
	req.Args = []string{"vet", "./"}
	return nil
}
```
