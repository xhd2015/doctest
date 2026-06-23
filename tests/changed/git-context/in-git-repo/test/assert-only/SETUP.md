# Scenario

**Feature**: changing one leaf `ASSERT.md` runs only that leaf

```
# only leaf_a ASSERT.md modified
changed leaf_a/ASSERT.md -> doctest test --changed -> 1 Run
```

## Steps

1. Create flat two-leaf tree and commit.
2. Modify `leaf_a/ASSERT.md` (unstaged).
3. Run `doctest test --changed`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	assertPath := filepath.Join(fx.TreeDir, "leaf_a", "ASSERT.md")
	content, err := os.ReadFile(assertPath)
	if err != nil {
		t.Fatal(err)
	}
	content = append(content, []byte("\n<!-- changed -->\n")...)
	if err := os.WriteFile(assertPath, content, 0644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = fx.RepoDir
	req.Args = []string{"test", fx.TreeDir, "--changed"}
	return nil
}
```