package core

import (
	"os"
	"path/filepath"
	"testing"
)

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
	desired := map[string]struct{}{
		"tree-a/keep.go": {},
		"go.mod":         {},
	}
	if err := PruneTreeScopeToDesired(gen, "tree-a", desired); err != nil {
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
	if err := PruneTreeScopeToDesired(gen, "tree-a", nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(f); err != nil {
		t.Fatalf("empty desired must not delete: %v", err)
	}
}
