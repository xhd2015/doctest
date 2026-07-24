package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureCleanGenRootSessionWipeOnce(t *testing.T) {
	gen := t.TempDir()
	orphan := filepath.Join(gen, "orphan.txt")
	if err := os.WriteFile(orphan, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := EnsureCleanGenRoot(gen, "sess-1", false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatalf("orphan should be wiped on first session visit, err=%v", err)
	}
	// recreate orphan; same session must not wipe
	if err := os.WriteFile(orphan, []byte("again"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := EnsureCleanGenRoot(gen, "sess-1", false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphan); err != nil {
		t.Fatalf("same session should not re-wipe: %v", err)
	}
	// new session wipes
	if err := EnsureCleanGenRoot(gen, "sess-2", false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatalf("new session should wipe orphan")
	}
}

func TestEnsureCleanGenRootNoSessionNoWipe(t *testing.T) {
	gen := t.TempDir()
	orphan := filepath.Join(gen, "keep.txt")
	if err := os.WriteFile(orphan, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := EnsureCleanGenRoot(gen, "", false); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphan); err != nil {
		t.Fatalf("library path without session must not wipe: %v", err)
	}
}

func TestEnsureCleanGenRootForceAWithoutSession(t *testing.T) {
	gen := t.TempDir()
	orphan := filepath.Join(gen, "gone.txt")
	if err := os.WriteFile(orphan, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := EnsureCleanGenRoot(gen, "", true); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Fatal("forceA without session should wipe")
	}
}
