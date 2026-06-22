# Scenario

**Feature**: one valid doctest tree and one invalid directory (missing DOCTEST.md)

```
# inspect test tree for structural issues
doctest vet <dir> -> walk tree -> report anti-patterns

# anti-patterns detected
embedded go block | go test shellout | assert without setup | skipped testdata
```

## Preconditions
- One valid doctest tree and one invalid directory (missing DOCTEST.md).

## Steps
1. Create one valid temp directory with a doctest tree.
2. Create a second temp directory without DOCTEST.md (invalid).
3. Run `doctest vet <valid-dir> <invalid-dir>`.

```go
import (
    "github.com/xhd2015/doctest/libdoc/testtree"
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	validDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(validDir, "DOCTEST.md"), []byte(testtree.VetDOCTEST()), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(validDir, "SETUP.md"), []byte("# Scenario\n\n**Feature**: minimal test setup\n\n\x60\x60\x60\n# minimal pipeline\nsystem -> run\n\x60\x60\x60\n\n## Setup\n"), 0644); err != nil {
		t.Fatal(err)
	}

	invalidDir := t.TempDir()

	req.Args = []string{"vet", validDir, invalidDir}
	return nil
}
```
