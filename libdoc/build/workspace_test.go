package build

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

// TestSharedGenRootAlwaysAssertReplace ensures multi-tree prepare under one
// external module always writes assert-mod replace (not only when a tree
// imports assert). Otherwise a tree without assert strips the replace and
// tidy resolves assert via github.com/xhd2015/doctest@vX → ambiguous session.
func TestSharedGenRootAlwaysAssertReplace(t *testing.T) {
	mod := t.TempDir()
	if err := os.WriteFile(filepath.Join(mod, "go.mod"), []byte("module example.com/assertfix\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}
	// Minimal trees: harness does not import github.com/xhd2015/doctest/assert.
	treeA := filepath.Join(mod, "tree-a")
	treeB := filepath.Join(mod, "tree-b")
	testtree.WriteMinimalRunnableTree(t, treeA, []testtree.LeafSpec{
		{Name: "a1", Steps: "a1", Expected: "ok"},
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
	if _, err := PrepareTree(treeA, opts); err != nil {
		t.Fatalf("PrepareTree A: %v", err)
	}
	if _, err := PrepareTree(treeB, opts); err != nil {
		t.Fatalf("PrepareTree B: %v", err)
	}
	goModData, err := os.ReadFile(filepath.Join(genDir, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	goMod := string(goModData)
	if !strings.Contains(goMod, "replace github.com/xhd2015/doctest/assert =>") {
		t.Fatalf("expected always-on assert replace after trees without assert import:\n%s", goMod)
	}
	if !strings.Contains(goMod, "replace github.com/xhd2015/doctest/session =>") {
		t.Fatalf("expected session replace:\n%s", goMod)
	}
	if !strings.Contains(goMod, "replace example.com/assertfix =>") {
		t.Fatalf("expected parent module replace:\n%s", goMod)
	}
	// Tidy must not hit ambiguous session (would if assert resolved via doctest@vX).
	if err := core.CondTidyGoMod(genDir); err != nil {
		t.Fatalf("CondTidyGoMod: %v", err)
	}
}

// TestPrepareTreeParallelSharedGenRoot stress-tests concurrent writeCases into
// one gen root (P3: no global writeCases lock; genModMu only on shared files).
func TestPrepareTreeParallelSharedGenRoot(t *testing.T) {
	mod := t.TempDir()
	if err := os.WriteFile(filepath.Join(mod, "go.mod"), []byte("module example.com/par\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}
	const n = 6
	trees := make([]string, n)
	for i := 0; i < n; i++ {
		trees[i] = filepath.Join(mod, "t"+string(rune('a'+i)))
		testtree.WriteMinimalRunnableTree(t, trees[i], []testtree.LeafSpec{
			{Name: "one", Steps: "s", Expected: "ok"},
			{Name: "two", Steps: "s2", Expected: "ok"},
		})
	}
	genDir := filepath.Join(t.TempDir(), "gen")
	opts := core.Options{
		GenDir:                genDir,
		Count:                 1,
		SuppressResultSummary: true,
		Stderr:                ioDiscard{},
	}
	preps := make([]TreePrep, n)
	errs := make([]error, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			p, err := PrepareTree(trees[i], opts)
			preps[i] = p
			errs[i] = err
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Fatalf("PrepareTree %d: %v", i, err)
		}
		if preps[i].Stats.Total != 2 {
			t.Fatalf("tree %d total=%d", i, preps[i].Stats.Total)
		}
	}
	var stdout bytes.Buffer
	runOpts := opts
	runOpts.Stdout = &stdout
	stats, err := RunWorkspace(preps, runOpts)
	if err != nil {
		t.Fatalf("RunWorkspace: %v\n%s", err, stdout.String())
	}
	if stats.Passed != n*2 || stats.Total != n*2 {
		t.Fatalf("want %d/%d, got %d/%d\n%s", n*2, n*2, stats.Passed, stats.Total, stdout.String())
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }
