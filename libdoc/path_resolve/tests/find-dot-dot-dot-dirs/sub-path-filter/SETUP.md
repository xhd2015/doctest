# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Steps
- Create a temp directory with `go.mod` and `subp/` containing `DOCTEST.md`.
- Also create `sibling/` without `DOCTEST.md`.
- Set `req.BasePath` to `subp/`.

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
	subp := filepath.Join(dir, "subp")
	if err := os.MkdirAll(subp, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subp, "DOCTEST.md"), []byte("# subp\n"), 0644); err != nil {
		t.Fatal(err)
	}
	sibling := filepath.Join(dir, "sibling")
	if err := os.MkdirAll(sibling, 0755); err != nil {
		t.Fatal(err)
	}
	req.BasePath = subp
	return nil
}
```
