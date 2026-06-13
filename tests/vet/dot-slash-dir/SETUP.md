## Preconditions
- A valid doctest tree exists in the current working directory.

## Steps
1. Create a minimal valid doctest tree in a temp directory.
2. Run `doctest vet ./` with WorkDir set to that directory.

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
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte("## Setup\n"), 0644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = dir
	req.Args = []string{"vet", "./"}
	return nil
}
```
