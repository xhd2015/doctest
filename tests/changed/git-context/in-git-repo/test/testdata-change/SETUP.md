# Scenario

**Feature**: changing a file under leaf `testdata/` runs that leaf

```
# testdata file modified
changed leaf_a/testdata/input.txt -> doctest test --changed -> 1 Run
```

## Steps

1. Create flat two-leaf tree with `leaf_a/testdata/input.txt` and commit.
2. Modify `testdata/input.txt`.
3. Run `doctest test --changed`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	fx := createFlatTwoLeafTree(t)
	testdataDir := filepath.Join(fx.TreeDir, "leaf_a", "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testdataDir, "input.txt"), []byte("baseline\n"), 0644); err != nil {
		t.Fatal(err)
	}
	gitAddCommitAll(t, fx.RepoDir, "add testdata")
	if err := os.WriteFile(filepath.Join(testdataDir, "input.txt"), []byte("modified\n"), 0644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = fx.RepoDir
	req.Args = []string{"test", fx.TreeDir, "--changed"}
	return nil
}
```