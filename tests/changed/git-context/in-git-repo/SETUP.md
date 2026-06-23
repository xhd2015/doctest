# Scenario

**Feature**: `--changed` filters doctest leaves by git working-tree changes

```
# ephemeral git repo with fixture tree
git init -> commit baseline -> modify paths -> doctest --changed

# map changed files to leaves
changed ASSERT.md -> one leaf | changed group SETUP.md -> descendant leaves
```

## Preconditions

- Git is available on PATH.
- Fixture trees use valid `DOCTEST.md` with `Request`, `Response`, and `Run`.

## Steps

1. Create an ephemeral git repository with a committed fixture tree.
2. Apply the leaf-specific file change (staged, unstaged, or untracked).
3. Run the doctest subcommand with `--changed` from the repo root.

## Context

- Helpers `initGitRepo`, `createFlatTwoLeafTree`, and `createSharedParentTwoLeafTree`
  build reproducible fixtures.
- Summary line parsing verifies how many leaves ran.

```go
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

const tick = "\u0060"

const leafSetupGo = "import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }"

const leafAssertGo = "import \"testing\"\n\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}"

const fixtureSetupGo = "import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }"

type gitFixture struct {
	RepoDir string
	TreeDir string
}

func goFence(n int) string {
	return strings.Repeat(tick, n)
}

func goBlock(code string) string {
	fence := goFence(3)
	return fence + "go\n" + code + "\n" + fence + "\n"
}

func scenarioHeader(feature, snippet string) string {
	fence := goFence(3)
	return fmt.Sprintf("# Scenario\n\n**Feature**: %s\n\n%s\n%s\n%s\n\n", feature, fence, snippet, fence)
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %s in %s: %v\n%s", strings.Join(args, " "), dir, err, string(out))
	}
}

func initGitRepo(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "changed@test.com")
	runGit(t, dir, "config", "user.name", "Changed Test")
}

func gitAddCommitAll(t *testing.T, dir, msg string) {
	t.Helper()
	runGit(t, dir, "add", "-A")
	runGit(t, dir, "commit", "-m", msg)
}

func writeLeaf(t *testing.T, root, name string) {
	t.Helper()
	leafDir := filepath.Join(root, name)
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		t.Fatal(err)
	}
	setup := scenarioHeader("fixture leaf", "# leaf setup\nleaf -> run") + "## Steps\n1. setup\n\n" + goBlock(leafSetupGo)
	assert := "## Expected\n- passes\n\n" + goBlock(leafAssertGo)
	if err := os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(setup), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(assert), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeRootTree(t *testing.T, root string, withRootSetup bool) {
	t.Helper()
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(testtree.MinimalRunGo())), 0644); err != nil {
		t.Fatal(err)
	}
	if withRootSetup {
		setup := scenarioHeader("fixture root", "# root setup\nroot -> shared context") + "## Steps\n1. root setup\n\n" + goBlock(fixtureSetupGo)
		if err := os.WriteFile(filepath.Join(root, "SETUP.md"), []byte(setup), 0644); err != nil {
			t.Fatal(err)
		}
	}
}

func createFlatTwoLeafTree(t *testing.T) gitFixture {
	t.Helper()
	repoDir := t.TempDir()
	treeDir := filepath.Join(repoDir, "tests", "fixture")
	writeRootTree(t, treeDir, true)
	writeLeaf(t, treeDir, "leaf_a")
	writeLeaf(t, treeDir, "leaf_b")
	initGitRepo(t, repoDir)
	gitAddCommitAll(t, repoDir, "baseline two leaves")
	return gitFixture{RepoDir: repoDir, TreeDir: treeDir}
}

func createSharedParentTwoLeafTree(t *testing.T) gitFixture {
	t.Helper()
	repoDir := t.TempDir()
	treeDir := filepath.Join(repoDir, "tests", "fixture")
	writeRootTree(t, treeDir, true)
	groupDir := filepath.Join(treeDir, "shared")
	if err := os.MkdirAll(groupDir, 0755); err != nil {
		t.Fatal(err)
	}
	groupSetup := scenarioHeader("shared parent", "# group setup\nshared -> both leaves") + "## Steps\n1. group setup\n\n" + goBlock(fixtureSetupGo)
	if err := os.WriteFile(filepath.Join(groupDir, "SETUP.md"), []byte(groupSetup), 0644); err != nil {
		t.Fatal(err)
	}
	writeLeaf(t, groupDir, "leaf_a")
	writeLeaf(t, groupDir, "leaf_b")
	initGitRepo(t, repoDir)
	gitAddCommitAll(t, repoDir, "baseline shared parent")
	return gitFixture{RepoDir: repoDir, TreeDir: treeDir}
}

func findInlineSummaryLine(stdout string) string {
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, " Run, ") && strings.Contains(line, " Pass, ") {
			return line
		}
	}
	return ""
}

func countGeneratedTestGoFiles(t *testing.T, genDir string) int {
	t.Helper()
	count := 0
	err := filepath.Walk(genDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "_test.go") {
			count++
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk gen dir: %v", err)
	}
	return count
}

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "CHANGED_TEST_GROUP=in-git-repo")
	return nil
}
```