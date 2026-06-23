# Scenario

**Feature**: only non-doctest file changes yield no-tests-changed warning

```
# README.md modified, no doctest files touched
changed README.md -> doctest test --changed -> warning, exit 0
```

## Steps

1. Create flat two-leaf tree and commit.
2. Modify `README.md` outside the test tree.
3. Run `doctest test --changed`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	if err := os.WriteFile(filepath.Join(fx.RepoDir, "README.md"), []byte("# changed\n"), 0644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = fx.RepoDir
	req.Args = []string{"test", fx.TreeDir, "--changed"}
	return nil
}
```