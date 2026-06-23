# Scenario

**Feature**: vet fails only on changed leaf with anti-pattern

```
# leaf_b SETUP.md changed to anti-pattern; leaf_a SETUP.md clean and unchanged
doctest vet --changed -> fail on leaf_b only
```

## Steps

1. Create two-leaf tree with clean SETUP files and commit.
2. Replace `leaf_b/SETUP.md` with embedded-Go anti-pattern (unstaged).
3. Run `doctest vet --changed`.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	fx := createVetFlatTwoLeafTree(t)
	setupPath := filepath.Join(fx.TreeDir, "leaf_b", "SETUP.md")
	if err := os.WriteFile(setupPath, readAntiPatternSetup(t), 0644); err != nil {
		t.Fatal(err)
	}
	req.WorkDir = fx.RepoDir
	req.Args = []string{"vet", fx.TreeDir, "--changed"}
	return nil
}
```