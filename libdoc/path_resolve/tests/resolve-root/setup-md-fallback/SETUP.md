## Steps
- Create a temp directory with only `SETUP.md` (no DOCTEST.md).
- Set `req.Input` to a sub-path within that directory.
- `ResolveRoot` should fall back to the directory containing `SETUP.md`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte("## Setup\n"), 0644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "sub", "leaf")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	req.Input = sub
	return nil
}
```
