package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWipeGenRoot_cleansExposeProductFiles(t *testing.T) {
	t.Parallel()
	genRoot := t.TempDir()
	t.Cleanup(func() { unregisterExposeGenRoot(genRoot) })
	product := t.TempDir()
	virt := filepath.Join(product, DoctestInternalExposeDir, "greet", "expose.go")
	if err := os.MkdirAll(filepath.Dir(virt), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(virt, []byte("package greet\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := recordExposeMaterialized(genRoot, virt); err != nil {
		t.Fatal(err)
	}
	if err := WipeGenRoot(genRoot); err != nil {
		t.Fatalf("WipeGenRoot: %v", err)
	}
	if _, err := os.Stat(virt); !os.IsNotExist(err) {
		t.Fatalf("product expose must be stripped before gen wipe: %v", err)
	}
	if _, err := os.Stat(filepath.Join(genRoot, ExposeMaterializedList)); !os.IsNotExist(err) {
		t.Fatal("list should be gone after successful wipe")
	}
}

func TestWipeGenRoot_failsClosedOnExposeCleanupError(t *testing.T) {
	t.Parallel()
	genRoot := t.TempDir()
	t.Cleanup(func() { unregisterExposeGenRoot(genRoot) })
	product := t.TempDir()
	virt := filepath.Join(product, DoctestInternalExposeDir, "greet", "expose.go")
	if err := os.MkdirAll(virt, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(virt, "keep.txt"), []byte("x\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := recordExposeMaterialized(genRoot, virt); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(genRoot, "keep-me.txt")
	if err := os.WriteFile(marker, []byte("stay\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := WipeGenRoot(genRoot); err == nil {
		t.Fatal("expected WipeGenRoot to fail when expose leftover cannot be removed")
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("gen root must not be wiped after cleanup failure: %v", err)
	}
	if _, err := os.Stat(filepath.Join(genRoot, ExposeMaterializedList)); err != nil {
		t.Fatalf("list must remain after failed wipe: %v", err)
	}
	if _, err := os.Stat(virt); err != nil {
		t.Fatalf("poison expose dir should remain: %v", err)
	}
}

func TestRootBookkeeping_exposeList(t *testing.T) {
	t.Parallel()
	if !rootBookkeeping(ExposeMaterializedList) {
		t.Fatal("expose materialized list must be gen-root bookkeeping")
	}
}

func TestGenBatchWipeOnce(t *testing.T) {
	gen := t.TempDir()
	orphan := filepath.Join(gen, "orphan.txt")
	if err := os.WriteFile(orphan, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	b := NewGenBatch()
	if err := b.WipeOnce(gen); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatal("first WipeOnce should remove orphan")
	}
	if err := os.WriteFile(orphan, []byte("again"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := b.WipeOnce(gen); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphan); err != nil {
		t.Fatalf("second WipeOnce same batch must not re-wipe: %v", err)
	}
}

func TestPruneTreeScopeDoesNotTouchSiblingTree(t *testing.T) {
	gen := t.TempDir()
	// tree-a keep + orphan, tree-b must survive prune of tree-a
	keepA := filepath.Join(gen, "tree-a", "keep.go")
	orphanA := filepath.Join(gen, "tree-a", "orphan.go")
	keepB := filepath.Join(gen, "tree-b", "keep.go")
	for _, p := range []string{keepA, orphanA, keepB} {
		if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("package x\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(gen, "go.mod"), []byte("module testcase\n"), 0644); err != nil {
		t.Fatal(err)
	}
	// Orphans must be in gen-manifest to be pruned (safe mode).
	if _, err := WriteIfChanged(gen, "tree-a/keep.go", []byte("package x\n")); err != nil {
		t.Fatal(err)
	}
	if _, err := WriteIfChanged(gen, "tree-a/orphan.go", []byte("package x\n")); err != nil {
		t.Fatal(err)
	}
	if _, err := WriteIfChanged(gen, "tree-b/keep.go", []byte("package x\n")); err != nil {
		t.Fatal(err)
	}
	desired := map[string]struct{}{
		"tree-a/keep.go": {},
		"go.mod":         {},
	}
	if err := PruneTreeScopeToDesired(gen, "tree-a", desired, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(keepA); err != nil {
		t.Fatalf("keepA: %v", err)
	}
	if _, err := os.Stat(orphanA); !os.IsNotExist(err) {
		t.Fatal("orphanA must be deleted")
	}
	if _, err := os.Stat(keepB); err != nil {
		t.Fatalf("sibling tree-b must not be touched: %v", err)
	}
	if _, err := os.Stat(filepath.Join(gen, "go.mod")); err != nil {
		t.Fatalf("go.mod must remain: %v", err)
	}
}

func TestNoteDesiredViaAttach(t *testing.T) {
	gen := t.TempDir()
	b := NewGenBatch()
	unlock := LockGenRootWrites(gen)
	defer unlock()
	b.Attach(gen)
	defer b.Detach(gen)
	NoteDesired(gen, "pkg/x.go")
	d := b.Desired(gen)
	if _, ok := d["pkg/x.go"]; !ok {
		t.Fatalf("expected pkg/x.go in desired, got %#v", d)
	}
}

func TestPruneEmptyDesiredNoOp(t *testing.T) {
	gen := t.TempDir()
	f := filepath.Join(gen, "tree-a", "x.txt")
	if err := os.MkdirAll(filepath.Dir(f), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(f, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := PruneTreeScopeToDesired(gen, "tree-a", nil, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(f); err != nil {
		t.Fatalf("empty desired must not delete: %v", err)
	}
}
