## Preconditions
- A temp Go module with a subdirectory containing DOCTEST.md and a sibling without.

## Steps
1. Create a temp directory with go.mod.
2. Create `subp/` with DOCTEST.md.
3. Create `sibling/` without DOCTEST.md.
4. Run `doctest vet ./subp/...` with WorkDir set to the temp directory.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644); err != nil {
		t.Fatal(err)
	}

	subp := filepath.Join(dir, "subp")
	if err := os.MkdirAll(subp, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subp, "DOCTEST.md"), []byte("# subp tests\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subp, "SETUP.md"), []byte("## Setup\n"), 0644); err != nil {
		t.Fatal(err)
	}

	sibling := filepath.Join(dir, "sibling")
	if err := os.MkdirAll(sibling, 0755); err != nil {
		t.Fatal(err)
	}

	req.WorkDir = dir
	req.Args = []string{"vet", "./subp/..."}
	return nil
}
```
