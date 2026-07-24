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

func TestPruneGenRootToDesired(t *testing.T) {
	gen := t.TempDir()
	keep := filepath.Join(gen, "keep", "a.go")
	if err := os.MkdirAll(filepath.Dir(keep), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(keep, []byte("package keep\n"), 0644); err != nil {
		t.Fatal(err)
	}
	orphan := filepath.Join(gen, "stale", "b.go")
	if err := os.MkdirAll(filepath.Dir(orphan), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(orphan, []byte("package stale\n"), 0644); err != nil {
		t.Fatal(err)
	}
	desired := map[string]struct{}{"keep/a.go": {}}
	if err := PruneGenRootToDesired(gen, desired); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatalf("keep must remain: %v", err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatal("orphan must be deleted")
	}
}

func TestNoteDesiredViaAttach(t *testing.T) {
	gen := t.TempDir()
	b := NewGenBatch()
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
	f := filepath.Join(gen, "x.txt")
	if err := os.WriteFile(f, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := PruneGenRootToDesired(gen, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(f); err != nil {
		t.Fatalf("empty desired must not delete: %v", err)
	}
}
