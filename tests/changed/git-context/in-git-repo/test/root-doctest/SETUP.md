# Scenario

**Feature**: changing root `DOCTEST.md` runs all leaves in the tree

```
# root DOCTEST.md modified
changed root DOCTEST.md -> doctest test --changed -> all leaves run
```

## Steps

1. Create flat two-leaf tree and commit.
2. Modify root `DOCTEST.md`.
3. Run `doctest test --changed`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	doctestPath := filepath.Join(fx.TreeDir, "DOCTEST.md")
	content, err := os.ReadFile(doctestPath)
	if err != nil {
		t.Fatal(err)
	}
	content = append(content, []byte("\n<!-- changed -->\n")...)
	if err := os.WriteFile(doctestPath, content, 0644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = fx.RepoDir
	req.Args = []string{"test", fx.TreeDir, "--changed"}
	return nil
}
```