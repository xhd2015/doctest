## Steps
- Create a temp directory with a `DOCTEST.md` file.
- Set `req.Input` to a sub-path within that directory.
- `ResolveRoot` should walk up and find the root.

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
	sub := filepath.Join(dir, "sub", "leaf")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	req.Input = sub
	return nil
}
```
