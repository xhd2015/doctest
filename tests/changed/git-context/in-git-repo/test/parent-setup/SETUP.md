# Scenario

**Feature**: changing a parent `SETUP.md` runs all descendant leaves

```
# shared/SETUP.md modified
changed group SETUP.md -> doctest test --changed -> 2 Run
```

## Steps

1. Create shared-parent two-leaf tree and commit.
2. Modify `shared/SETUP.md`.
3. Run `doctest test --changed`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	fx := createSharedParentTwoLeafTree(t)
	setupPath := filepath.Join(fx.TreeDir, "shared", "SETUP.md")
	content, err := os.ReadFile(setupPath)
	if err != nil {
		t.Fatal(err)
	}
	content = append(content, []byte("\n<!-- changed -->\n")...)
	if err := os.WriteFile(setupPath, content, 0644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = fx.RepoDir
	req.Args = []string{"test", fx.TreeDir, "--changed"}
	return nil
}
```