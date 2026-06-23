# Scenario

**Feature**: `doctest build --changed` compiles only the changed leaf

```
# only leaf_a ASSERT.md modified
changed leaf_a/ASSERT.md -> doctest build --changed -> 1 generated test file
```

## Steps

1. Create flat two-leaf tree and commit.
2. Modify `leaf_a/ASSERT.md`.
3. Run `doctest build --changed --gen-dir`.

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
	content = append(content, []byte("\n<!-- build-changed -->\n")...)
	if err := os.WriteFile(assertPath, content, 0644); err != nil {
		t.Fatal(err)
	}
	genDir := filepath.Join(fx.RepoDir, "gen")
	req.WorkDir = fx.RepoDir
	req.Args = []string{"build", fx.TreeDir, "--changed", "--gen-dir", genDir}
	req.Env = append(req.Env, "CHANGED_GEN_DIR="+genDir)
	return nil
}
```