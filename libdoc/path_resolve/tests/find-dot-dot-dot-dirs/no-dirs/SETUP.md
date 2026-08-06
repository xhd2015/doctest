# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Steps
- Create a temp directory with `go.mod` but no `DOCTEST.md` anywhere.
- Set `req.BasePath` to the temp dir.

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
	req.BasePath = dir
	return nil
}
```
