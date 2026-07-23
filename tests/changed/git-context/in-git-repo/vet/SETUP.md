# Scenario

**Feature**: vet `--changed` selects doctest markdown via `ChangedDoctestMarkdownFiles`

```
# vet selection is path-based, not leaf-filter
ChangedDoctestMarkdownFiles(tree, gitRoot, changed) -> absolute markdown paths
```

## Preconditions

- Fixture tree under synthetic git root.
- Vet validates only returned markdown files; root is omitted when unchanged.

## Steps

1. Create fixture tree (valid or intentionally invalid root for skip-root).
2. Set `Policy=vet-md` and synthetic changed paths.
3. Assert relative markdown path list.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

func readAntiPatternSetup(t *testing.T, d *session.Doctest) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(d.DOCTEST_ROOT, "git-context", "in-git-repo", "vet", "fixture_anti_pattern_setup.md.txt"))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func writeInvalidRootDOCTEST(t *testing.T, root string) {
	t.Helper()
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	fence := strings.Repeat("\u0060", 3)
	// Missing ## Version — invalid for full vet, but skip-root must not select it when unchanged.
	content := "# Tests\n\n" + testtree.MinimalDSN + "\n" + fence + "go\n" + testtree.MinimalRunGo() + "\n" + fence + "\n"
	if err := os.WriteFile(filepath.Join(root, "DOCTEST.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func createVetFlatTwoLeafTree(t *testing.T) policyFixture {
	t.Helper()
	repoDir := t.TempDir()
	treeDir := filepath.Join(repoDir, "tests", "fixture")
	writeRootTree(t, treeDir, true)
	writeLeaf(t, treeDir, "leaf_a")
	writeLeaf(t, treeDir, "leaf_b")
	return policyFixture{RepoDir: repoDir, TreeDir: treeDir}
}

func createVetSkipRootTree(t *testing.T) policyFixture {
	t.Helper()
	repoDir := t.TempDir()
	treeDir := filepath.Join(repoDir, "tests", "fixture")
	writeInvalidRootDOCTEST(t, treeDir)
	writeLeaf(t, treeDir, "leaf_a")
	return policyFixture{RepoDir: repoDir, TreeDir: treeDir}
}

func Setup(t *testing.T, req *Request) error {
	req.Policy = PolicyVetMD
	return nil
}
```
