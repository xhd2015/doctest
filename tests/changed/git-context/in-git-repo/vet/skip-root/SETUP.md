# Scenario

**Feature**: vet skips unchanged root `DOCTEST.md` even when invalid

```
# root DOCTEST.md missing Version (invalid but committed, unchanged)
# only leaf_a ASSERT.md changed
doctest vet --changed -> exit 0 (root not validated)
```

## Steps

1. Create tree with invalid root `DOCTEST.md` (no `## Version`) and commit.
2. Modify only `leaf_a/ASSERT.md`.
3. Run `doctest vet --changed`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	fx := createVetSkipRootTree(t)
	assertPath := filepath.Join(fx.TreeDir, "leaf_a", "ASSERT.md")
	content, err := os.ReadFile(assertPath)
	if err != nil {
		t.Fatal(err)
	}
	content = append(content, []byte("\n<!-- vet skip root -->\n")...)
	if err := os.WriteFile(assertPath, content, 0644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = fx.RepoDir
	req.Args = []string{"vet", fx.TreeDir, "--changed"}
	return nil
}
```