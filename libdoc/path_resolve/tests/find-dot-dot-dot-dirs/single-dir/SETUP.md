## Steps
- Create a temp directory with `go.mod` and a single `DOCTEST.md`.
- Set `req.BasePath = "."`.
- Set `req.Input` as the temp dir (used as cwd via WorkDir equivalent).

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte("# single\n"), 0644); err != nil {
		t.Fatal(err)
	}
	req.BasePath = dir
	return nil
}
```
