package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

func TestFilterByChangedFilesLeafAssert(t *testing.T) {
	_, treeDir := changedFixtureRepo(t)
	assertPath := filepath.Join(treeDir, "leaf_a", "ASSERT.md")
	content, err := os.ReadFile(assertPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(assertPath, append(content, []byte("\n<!-- changed -->\n")...), 0644); err != nil {
		t.Fatal(err)
	}

	gitRoot, changed, err := ChangedGitFiles(treeDir)
	if err != nil {
		t.Fatal(err)
	}
	cases, err := DiscoverTreeCases(treeDir)
	if err != nil {
		t.Fatal(err)
	}
	filtered := FilterByChangedFiles(cases, treeDir, gitRoot, changed)
	if len(filtered) != 1 || filtered[0].Path != "leaf_a" {
		t.Fatalf("filtered = %#v, want [leaf_a]", filtered)
	}
}

func changedFixtureRepo(t *testing.T) (repoDir, treeDir string) {
	t.Helper()
	repoDir = t.TempDir()
	treeDir = filepath.Join(repoDir, "tests", "fixture")
	if err := os.MkdirAll(treeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(treeDir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(testtree.MinimalRunGo())), 0644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"leaf_a", "leaf_b"} {
		leaf := filepath.Join(treeDir, name)
		if err := os.MkdirAll(leaf, 0755); err != nil {
			t.Fatal(err)
		}
		setup := "# Scenario\n\n```go\nimport \"testing\"\nfunc Setup(t *testing.T, req *Request) error { return nil }\n```\n"
		assert := "```go\nimport \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n```\n"
		if err := os.WriteFile(filepath.Join(leaf, "SETUP.md"), []byte(setup), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(leaf, "ASSERT.md"), []byte(assert), 0644); err != nil {
			t.Fatal(err)
		}
	}
	runGit(t, repoDir, "init")
	runGit(t, repoDir, "config", "user.email", "changed@test.com")
	runGit(t, repoDir, "config", "user.name", "Changed Test")
	runGit(t, repoDir, "add", "-A")
	runGit(t, repoDir, "commit", "-m", "baseline")
	return repoDir, treeDir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}