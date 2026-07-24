package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

func TestFilterByChangedFilesIgnoresSiblingStrayFiles(t *testing.T) {
	_, treeDir := changedFixtureRepo(t)
	assertPath := filepath.Join(treeDir, "leaf_a", "ASSERT.md")
	content, err := os.ReadFile(assertPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(assertPath, append(content, []byte("\n<!-- changed -->\n")...), 0644); err != nil {
		t.Fatal(err)
	}
	strayPath := filepath.Join(treeDir, "leaf_b", "stray.go")
	if err := os.WriteFile(strayPath, []byte("package leaf_b\n"), 0644); err != nil {
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
		t.Fatalf("filtered = %#v, want [leaf_a] (sibling stray file must not widen runs)", filtered)
	}
}

func TestShouldAnnounceChangedRun(t *testing.T) {
	zero := ChangedRunInfo{ChangedCount: 0}
	changed := ChangedRunInfo{ChangedCount: 1}
	if ShouldAnnounceChangedRun(zero, false) {
		t.Fatal("zero-change tree should be silent without -v")
	}
	if !ShouldAnnounceChangedRun(zero, true) {
		t.Fatal("zero-change tree should announce with -v")
	}
	if !ShouldAnnounceChangedRun(changed, false) {
		t.Fatal("changed tree should announce without -v")
	}
}

func TestFormatDoctestAnnouncementChangedZero(t *testing.T) {
	line := FormatDoctestAnnouncement("go-pkgs/cmd/http-proxy-flex/tests", ChangedRunInfo{
		TotalInTree:  13,
		ChangedCount: 0,
	}, true, 0)
	want := "doctest: go-pkgs/cmd/http-proxy-flex/tests (13 tests, --changed: 0 tests)"
	if line != want {
		t.Fatalf("got %q, want %q", line, want)
	}
}

func TestFormatDoctestAnnouncementChangedWidened(t *testing.T) {
	line := FormatDoctestAnnouncement("go-pkgs/cmd/http-proxy-flex/tests", ChangedRunInfo{
		TotalInTree:  13,
		ChangedCount: 12,
		Detail:       "2 leaves + 1 DOCTEST.md affecting 10 other tests",
	}, true, 0)
	want := "doctest: go-pkgs/cmd/http-proxy-flex/tests (13 tests, --changed: 12 tests, 2 leaves + 1 DOCTEST.md affecting 10 other tests)"
	if line != want {
		t.Fatalf("got %q, want %q", line, want)
	}
}

func TestChangedRunInfoForTreeLeafAssertOnly(t *testing.T) {
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
	info := ChangedRunInfoForTree(cases, treeDir, gitRoot, changed)
	if info.TotalInTree != 2 || info.ChangedCount != 1 {
		t.Fatalf("info = %#v, want TotalInTree=2 ChangedCount=1", info)
	}
	if info.Detail != "1 leaf" {
		t.Fatalf("detail = %q, want %q", info.Detail, "1 leaf")
	}
}

func TestChangedRunInfoForTreeRootDoctest(t *testing.T) {
	_, treeDir := changedFixtureRepo(t)
	doctestPath := filepath.Join(treeDir, "DOCTEST.md")
	content, err := os.ReadFile(doctestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(doctestPath, append(content, []byte("\n<!-- changed -->\n")...), 0644); err != nil {
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
	info := ChangedRunInfoForTree(cases, treeDir, gitRoot, changed)
	if info.TotalInTree != 2 || info.ChangedCount != 2 {
		t.Fatalf("info = %#v, want TotalInTree=2 ChangedCount=2", info)
	}
	if info.Detail != "1 DOCTEST.md affecting 2 other tests" {
		t.Fatalf("detail = %q, want %q", info.Detail, "1 DOCTEST.md affecting 2 other tests")
	}
}

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
		setup := "# Scenario\n\n```go\nimport \"testing\"\nfunc Setup(t *testing.T, d *session.Doctest, req *Request) error { return nil }\n```\n"
		assert := "```go\nimport \"testing\"\nfunc Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {}\n```\n"
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