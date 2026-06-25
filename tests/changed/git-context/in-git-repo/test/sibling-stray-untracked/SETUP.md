# Scenario

**Feature**: unrelated untracked files in a sibling leaf must not widen `--changed` runs

```
# only leaf_a ASSERT.md modified; leaf_b has unrelated untracked file
changed leaf_a/ASSERT.md + untracked leaf_b/stray.go -> doctest test --changed ./... -> 1 Run
```

## Steps

1. Create flat two-leaf tree and commit.
2. Modify `leaf_a/ASSERT.md` (unstaged).
3. Add an unrelated untracked file under `leaf_b/`.
4. Run `doctest test --changed ./...` from the repo root.

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
	strayPath := filepath.Join(fx.TreeDir, "leaf_b", "stray.go")
	if err := os.WriteFile(strayPath, []byte("package leaf_b\n"), 0644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = fx.RepoDir
	req.Args = []string{"test", "--changed", "./..."}
	return nil
}
```