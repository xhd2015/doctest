package build

import (
	"bytes"
	"fmt"
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
	batch := core.NewGenBatch()
	opts := core.Options{
		GenDir:                genDir,
		GenBatch:              batch,
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
	batch := core.NewGenBatch()
	opts := core.Options{
		GenDir:                genDir,
		GenBatch:              batch,
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
	if err := core.CondTidyGoMod(genDir, ""); err != nil {
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
	batch := core.NewGenBatch()
	opts := core.Options{
		GenDir:                genDir,
		GenBatch:              batch,
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

// TestPrepareTree_siblingKindBSurvivesFailedPrepare is the generate-only
// contract: a failed PrepareTree must not strip Kind B files another tree
// wrote on the shared gen root. RunWorkspace then strips them.
func TestPrepareTree_siblingKindBSurvivesFailedPrepare(t *testing.T) {
	mod := t.TempDir()
	if err := os.WriteFile(filepath.Join(mod, "go.mod"), []byte("module example.com/app\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}
	greetDir := filepath.Join(mod, "internal", "greet")
	if err := os.MkdirAll(greetDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(greetDir, "greet.go"), []byte("package greet\n\nfunc Hello(name string) string { return \"hello \" + name }\n"), 0644); err != nil {
		t.Fatal(err)
	}

	treeA := filepath.Join(mod, "tree-a")
	treeB := filepath.Join(mod, "tree-b")
	writeTreeImportingInternal(t, treeA, "example.com/app/internal/greet", "greet")
	writeTreeImportingInternal(t, treeB, "example.com/app/internal/missing", "missing")

	genDir := filepath.Join(t.TempDir(), "gen")
	opts := core.Options{
		GenDir:                genDir,
		GenBatch:              core.NewGenBatch(),
		Count:                 1,
		SuppressResultSummary: true,
		Stderr:                ioDiscard{},
	}
	prepA, err := PrepareTree(treeA, opts)
	if err != nil {
		t.Fatalf("PrepareTree A: %v", err)
	}
	virt := filepath.Join(mod, core.DoctestInternalExposeDir, "greet", "expose.go")
	if _, err := os.Stat(virt); err != nil {
		t.Fatalf("A should materialize Kind B expose: %v", err)
	}

	prepB, err := PrepareTree(treeB, opts)
	if err == nil {
		t.Fatal("PrepareTree B: expected generate error for missing internal")
	}
	if prepB.GenRoot == "" {
		t.Fatal("failed PrepareTree must return GenRoot for leftover cleanup")
	}
	if filepath.Clean(prepB.GenRoot) != filepath.Clean(prepA.GenRoot) {
		t.Fatalf("expected shared gen root, got A=%q B=%q", prepA.GenRoot, prepB.GenRoot)
	}
	if _, err := os.Stat(virt); err != nil {
		t.Fatalf("sibling generate failure must not strip A's expose.go: %v", err)
	}

	var stdout bytes.Buffer
	runOpts := opts
	runOpts.Stdout = &stdout
	if _, err := RunWorkspace([]TreePrep{prepA}, runOpts); err != nil {
		t.Fatalf("RunWorkspace A: %v\n%s", err, stdout.String())
	}
	if _, err := os.Stat(virt); !os.IsNotExist(err) {
		t.Fatalf("expose.go should be stripped after RunWorkspace: %v", err)
	}
}

func writeTreeImportingInternal(t *testing.T, root, internalImp, pkgIdent string) {
	t.Helper()
	testtree.WriteMinimalRunnableTree(t, root, []testtree.LeafSpec{{Name: "one", Steps: "s", Expected: "ok"}})
	runGo := `import (
	"testing"

	"` + internalImp + `"
	"github.com/xhd2015/doctest/session"
)

type Request struct{}
type Response struct{}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	_ = ` + pkgIdent + `.Hello
	return &Response{}, nil
}`
	testtree.WriteFile(t, root, "DOCTEST.md", testtree.MinimalDOCTEST(runGo))
}

// TestCleanupKindBForPreps_failedPrepareGenRoot covers the all-fail session:
// no RunWorkspace, but the failed prep still names the shared gen root.
func TestCleanupKindBForPreps_failedPrepareGenRoot(t *testing.T) {
	mod := t.TempDir()
	if err := os.WriteFile(filepath.Join(mod, "go.mod"), []byte("module example.com/app\n\ngo 1.22\n"), 0644); err != nil {
		t.Fatal(err)
	}
	greetDir := filepath.Join(mod, "internal", "greet")
	if err := os.MkdirAll(greetDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(greetDir, "greet.go"), []byte("package greet\n\nfunc Hello(name string) string { return \"hi\" }\n"), 0644); err != nil {
		t.Fatal(err)
	}
	treeA := filepath.Join(mod, "tree-a")
	treeB := filepath.Join(mod, "tree-b")
	writeTreeImportingInternal(t, treeA, "example.com/app/internal/greet", "greet")
	writeTreeImportingInternal(t, treeB, "example.com/app/internal/missing", "missing")
	opts := core.Options{
		GenDir:                filepath.Join(t.TempDir(), "gen"),
		GenBatch:              core.NewGenBatch(),
		Count:                 1,
		SuppressResultSummary: true,
		Stderr:                ioDiscard{},
	}
	if _, err := PrepareTree(treeA, opts); err != nil {
		t.Fatalf("PrepareTree A: %v", err)
	}
	prepB, err := PrepareTree(treeB, opts)
	if err == nil {
		t.Fatal("expected B generate error")
	}
	virt := filepath.Join(mod, core.DoctestInternalExposeDir, "greet", "expose.go")
	if _, err := os.Stat(virt); err != nil {
		t.Fatalf("expose should remain after failed B: %v", err)
	}
	if err := CleanupKindBForPreps([]TreePrep{prepB}); err != nil {
		t.Fatalf("leftover cleanup: %v", err)
	}
	if _, err := os.Stat(virt); !os.IsNotExist(err) {
		t.Fatalf("failed prep GenRoot should be enough to strip Kind B: %v", err)
	}
}

// TestRunWorkspace_cleansKindBOnEarlyError covers the leak window where
// PrepareTree wrote product expose files but RunWorkspace returns before
// finishWorkspaceGoTest (e.g. non-unified prep).
func TestRunWorkspace_cleansKindBOnEarlyError(t *testing.T) {
	genRoot := t.TempDir()
	product := t.TempDir()
	virt := filepath.Join(product, core.DoctestInternalExposeDir, "greet", "expose.go")
	if err := os.MkdirAll(filepath.Dir(virt), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(virt, []byte("package greet\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(genRoot, core.KindBMaterializedList), []byte(virt+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := RunWorkspace([]TreePrep{{
		AbsRoot: product,
		Unified: false,
		GenRoot: genRoot,
		Stats:   TestRunStats{Total: 1},
	}}, core.Options{SuppressResultSummary: true, Stderr: ioDiscard{}})
	if err == nil {
		t.Fatal("expected error for non-unified prep")
	}
	if _, err := os.Stat(virt); !os.IsNotExist(err) {
		t.Fatalf("expose.go should be stripped on RunWorkspace error: %v", err)
	}
}

func TestJoinKindBCleanupErr(t *testing.T) {
	if err := joinKindBCleanupErr(nil, nil); err != nil {
		t.Fatalf("nil+nil: %v", err)
	}
	run := fmt.Errorf("workspace: not unified")
	if err := joinKindBCleanupErr(run, nil); err != run {
		t.Fatalf("nil cleanup should keep run err, got %v", err)
	}
	clean := fmt.Errorf("remove expose.go: is a directory")
	got := joinKindBCleanupErr(nil, clean)
	if got == nil || !strings.Contains(got.Error(), "kind B cleanup") || !strings.Contains(got.Error(), clean.Error()) {
		t.Fatalf("cleanup-only: %v", got)
	}
	mixed := joinKindBCleanupErr(run, clean)
	if mixed == nil || !strings.Contains(mixed.Error(), "not unified") || !strings.Contains(mixed.Error(), "kind B cleanup") {
		t.Fatalf("mixed: %v", mixed)
	}
}

func TestGenerateContextClose_surfacesKindBCleanupError(t *testing.T) {
	genRoot := t.TempDir()
	product := t.TempDir()
	virt := filepath.Join(product, core.DoctestInternalExposeDir, "greet", "expose.go")
	if err := os.MkdirAll(virt, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(virt, "keep.txt"), []byte("x\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(genRoot, core.KindBMaterializedList), []byte(virt+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	ctx := &generateContext{genRoot: genRoot, w: ioDiscard{}}
	err := ctx.Close()
	if err == nil {
		t.Fatal("expected cleanup error when expose.go is a non-empty dir")
	}
	if _, statErr := os.Stat(virt); statErr != nil {
		t.Fatalf("poison dir should remain: %v", statErr)
	}
}

func TestGenerateContextClose_cleansKindB(t *testing.T) {
	genRoot := t.TempDir()
	product := t.TempDir()
	virt := filepath.Join(product, core.DoctestInternalExposeDir, "greet", "expose.go")
	if err := os.MkdirAll(filepath.Dir(virt), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(virt, []byte("package greet\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(genRoot, core.KindBMaterializedList), []byte(virt+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	ctx := &generateContext{genRoot: genRoot, w: ioDiscard{}}
	if err := ctx.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if _, err := os.Stat(virt); !os.IsNotExist(err) {
		t.Fatalf("expose.go should be stripped: %v", err)
	}
}

func TestGenerateContextClose_generateOnlyLeavesKindB(t *testing.T) {
	genRoot := t.TempDir()
	product := t.TempDir()
	virt := filepath.Join(product, core.DoctestInternalExposeDir, "greet", "expose.go")
	if err := os.MkdirAll(filepath.Dir(virt), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(virt, []byte("package greet\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(genRoot, core.KindBMaterializedList), []byte(virt+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	ctx := &generateContext{genRoot: genRoot, generateOnly: true, w: ioDiscard{}}
	if err := ctx.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if _, err := os.Stat(virt); err != nil {
		t.Fatalf("generate-only must leave expose.go: %v", err)
	}
	if core.KindBInterruptExitEnabled() {
		t.Fatal("generate-only Close must not leave CLI os.Exit armed")
	}
}

func TestRunWorkspace_surfacesKindBCleanupError(t *testing.T) {
	genRoot := t.TempDir()
	product := t.TempDir()
	virt := filepath.Join(product, core.DoctestInternalExposeDir, "greet", "expose.go")
	if err := os.MkdirAll(virt, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(virt, "keep.txt"), []byte("x\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(genRoot, core.KindBMaterializedList), []byte(virt+"\n"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := RunWorkspace([]TreePrep{{
		AbsRoot: product,
		Unified: false,
		GenRoot: genRoot,
		Stats:   TestRunStats{Total: 1},
	}}, core.Options{SuppressResultSummary: true, Stderr: ioDiscard{}})
	if err == nil {
		t.Fatal("expected combined workspace + cleanup error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "not unified") {
		t.Fatalf("want original workspace error, got %v", err)
	}
	if !strings.Contains(msg, "kind B cleanup") {
		t.Fatalf("want cleanup error surfaced, got %v", err)
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }
