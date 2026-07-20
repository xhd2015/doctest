package build

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

func TestPrepareAndRunWorkspaceTwoTrees(t *testing.T) {
	// Module-like layout: two DOCTEST roots under one temp "module".
	mod := t.TempDir()
	if err := os.WriteFile(filepath.Join(mod, "go.mod"), []byte("module example.com/ws\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}
	treeA := filepath.Join(mod, "tree-a")
	treeB := filepath.Join(mod, "tree-b")
	testtree.WriteMinimalRunnableTree(t, treeA, []testtree.LeafSpec{
		{Name: "a1", Steps: "a1", Expected: "ok"},
		{Name: "a2", Steps: "a2", Expected: "ok"},
	})
	testtree.WriteMinimalRunnableTree(t, treeB, []testtree.LeafSpec{
		{Name: "b1", Steps: "b1", Expected: "ok"},
	})

	genDir := filepath.Join(t.TempDir(), "gen")
	opts := core.Options{
		GenDir:                genDir,
		Count:                 1,
		SuppressResultSummary: true,
		Stderr:                ioDiscard{},
	}

	prepA, err := PrepareTree(treeA, opts)
	if err != nil {
		t.Fatalf("PrepareTree A: %v", err)
	}
	if !prepA.Unified || prepA.Stats.Total != 2 {
		t.Fatalf("prep A: unified=%v total=%d", prepA.Unified, prepA.Stats.Total)
	}
	prepB, err := PrepareTree(treeB, opts)
	if err != nil {
		t.Fatalf("PrepareTree B: %v", err)
	}
	if prepB.Stats.Total != 1 {
		t.Fatalf("prep B total=%d", prepB.Stats.Total)
	}
	if filepath.Clean(prepA.GenRoot) != filepath.Clean(prepB.GenRoot) {
		t.Fatalf("expected shared gen root, got %q vs %q", prepA.GenRoot, prepB.GenRoot)
	}

	// Per-tree wreg written.
	for _, tr := range []string{prepA.TreeRel, prepB.TreeRel} {
		if _, err := os.Stat(filepath.Join(genDir, tr, core.TreeWregDirName, "wreg.go")); err != nil {
			t.Fatalf("wreg for %s: %v", tr, err)
		}
	}

	var stdout bytes.Buffer
	runOpts := opts
	runOpts.Stdout = &stdout
	stats, err := RunWorkspace([]TreePrep{prepA, prepB}, runOpts)
	if err != nil {
		t.Fatalf("RunWorkspace: %v\nstdout:\n%s", err, stdout.String())
	}
	if stats.Passed != 3 || stats.Total != 3 {
		t.Fatalf("want 3/3 pass, got pass=%d total=%d\n%s", stats.Passed, stats.Total, stdout.String())
	}
	// Three leaf dots (workspace nested counting).
	out := stdout.String()
	summaryIdx := strings.Index(out, "  (")
	if summaryIdx < 0 {
		t.Fatalf("missing summary:\n%s", out)
	}
	dots := strings.Count(out[:summaryIdx], ".")
	if dots != 3 {
		t.Fatalf("expected 3 progress dots before summary, got %d:\n%s", dots, out)
	}
	if _, err := os.Stat(filepath.Join(genDir, core.WorkspaceDirName, core.WorkspaceSuiteDirName, "suite_test.go")); err != nil {
		t.Fatalf("workspace suite missing: %v", err)
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }
