# Scenario

**Feature**: `doctest vet --changed` validates only changed doctest files

```
# vet changed files only
doctest vet --changed <tree> -> walk changed paths -> skip unchanged siblings and root
```

## Preconditions

- Fixture tree lives inside an initialized git repository.

## Steps

1. Create and commit a baseline fixture tree.
2. Apply leaf-specific changes to doctest markdown files.
3. Run `doctest vet <tree> --changed`.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

func readAntiPatternSetup(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(DOCTEST_ROOT, "git-context", "in-git-repo", "vet", "fixture_anti_pattern_setup.md.txt"))
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
	content := "# Tests\n\n" + testtree.MinimalDSN + "\n" + fence + "go\n" + testtree.MinimalRunGo() + "\n" + fence + "\n"
	if err := os.WriteFile(filepath.Join(root, "DOCTEST.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func createVetFlatTwoLeafTree(t *testing.T) gitFixture {
	t.Helper()
	repoDir := t.TempDir()
	treeDir := filepath.Join(repoDir, "tests", "fixture")
	writeRootTree(t, treeDir, true)
	writeLeaf(t, treeDir, "leaf_a")
	writeLeaf(t, treeDir, "leaf_b")
	initGitRepo(t, repoDir)
	gitAddCommitAll(t, repoDir, "vet baseline")
	return gitFixture{RepoDir: repoDir, TreeDir: treeDir}
}

func createVetSkipRootTree(t *testing.T) gitFixture {
	t.Helper()
	repoDir := t.TempDir()
	treeDir := filepath.Join(repoDir, "tests", "fixture")
	writeInvalidRootDOCTEST(t, treeDir)
	writeLeaf(t, treeDir, "leaf_a")
	initGitRepo(t, repoDir)
	gitAddCommitAll(t, repoDir, "invalid root committed")
	return gitFixture{RepoDir: repoDir, TreeDir: treeDir}
}

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "CHANGED_SUBCMD=vet")
	return nil
}
```