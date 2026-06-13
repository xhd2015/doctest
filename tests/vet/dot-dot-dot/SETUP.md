## Preconditions
- A temp Go module exists with multiple subdirectories containing DOCTEST.md.

## Steps
1. Create a temp directory with go.mod.
2. Create subdirectories `sub-a/` and `sub-b/` each with DOCTEST.md.
3. Run `doctest vet ./...` with WorkDir set to the temp directory.

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

	for _, name := range []string{"sub-a", "sub-b"} {
		sub := filepath.Join(dir, name)
		if err := os.MkdirAll(sub, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(sub, "DOCTEST.md"), []byte("# tests "+name+"\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(sub, "SETUP.md"), []byte("## Setup\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	req.WorkDir = dir
	req.Args = []string{"vet", "./..."}
	return nil
}
```
