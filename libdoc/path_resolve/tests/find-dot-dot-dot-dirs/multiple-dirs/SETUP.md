## Steps
- Create a temp directory with `go.mod` and three subdirectories each containing `DOCTEST.md`.
- Set `req.BasePath` to the temp dir.

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
	for _, name := range []string{"sub-a", "sub-b", "sub-c"} {
		sub := filepath.Join(dir, name)
		if err := os.MkdirAll(sub, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(sub, "DOCTEST.md"), []byte("# "+name+"\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	req.BasePath = dir
	return nil
}
```
