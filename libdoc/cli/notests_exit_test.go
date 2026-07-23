package cli

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/path_resolve"
	"github.com/xhd2015/doctest/libdoc/runner"
)

// Empty ./... under a module child with no doctests must soft-exit (nil error,
// "no tests" on stderr). Regression: FindDotDotDotDirs used to return a plain
// errors.New("no tests") that failed errors.Is(…, ErrNoTestsFound), so the
// process exited 1 — CI saw ExitCode 1 while some local paths appeared soft.
func TestNoTestsDotDotDotSoftExit(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module x\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}
	// Same layout as tests/test/dotdotdot/git-boundary/same-repo-child-dir-no-tests:
	// parent has a tree; child does not. ./... from child must not walk up.
	parent := filepath.Join(root, "parent_test", "simple")
	if err := os.MkdirAll(parent, 0755); err != nil {
		t.Fatal(err)
	}
	// Minimal DOCTEST so a buggy walk-up would find tests.
	if err := os.WriteFile(filepath.Join(root, "parent_test", "DOCTEST.md"), []byte("# T\n## Version\n0.0.2\n"), 0644); err != nil {
		t.Fatal(err)
	}
	child := filepath.Join(root, "child")
	if err := os.MkdirAll(child, 0755); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("git", "-C", root, "init").CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
	if err := os.Chdir(child); err != nil {
		t.Fatal(err)
	}

	_, err = path_resolve.FindDotDotDotDirs(".")
	if !errors.Is(err, path_resolve.ErrNoTestsFound) {
		t.Fatalf("FindDotDotDotDirs: want ErrNoTestsFound, got %v (isSentinel=%v)", err, errors.Is(err, path_resolve.ErrNoTestsFound))
	}

	err = runner.Test([]string{"./" + "..."})
	if !errors.Is(err, runner.ErrNoTestsFound) {
		t.Fatalf("runner.Test: want ErrNoTestsFound, got %v", err)
	}

	var out, errb bytes.Buffer
	err = RunWithWriters(&out, &errb, []string{"test", "./..."})
	if err != nil {
		t.Fatalf("RunWithWriters: want soft nil, got %v (stderr=%q)", err, errb.String())
	}
	if !bytes.Contains(errb.Bytes(), []byte("no tests")) {
		t.Fatalf("stderr must contain \"no tests\", got %q", errb.String())
	}
	if bytes.Contains(errb.Bytes(), []byte("parent_test")) {
		t.Fatalf("./... from child must not discover parent_test:\n%s", errb.String())
	}
}
