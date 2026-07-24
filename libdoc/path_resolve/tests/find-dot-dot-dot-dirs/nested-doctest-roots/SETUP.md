## Steps
- Create a temp module with a parent `DOCTEST.md` and a nested `mapping-gen/DOCTEST.md`.
- Set `req.BasePath` to `"."` and run `FindDotDotDotDirs` from that module root.

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
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte("# parent\n"), 0644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(dir, "mapping-gen")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "DOCTEST.md"), []byte("# nested\n"), 0644); err != nil {
		t.Fatal(err)
	}
	req.BasePath = dir
	return nil
}
```