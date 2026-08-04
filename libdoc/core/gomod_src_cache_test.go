package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteGoModSrcFingerprintWarmHit(t *testing.T) {
	modRoot := t.TempDir()
	genDir := filepath.Join(t.TempDir(), "gen")
	writeTreeFile(t, modRoot, "go.mod", "module example.com/a\n\ngo 1.21\n")
	writeTreeFile(t, modRoot, "go.sum", "example.com/a v1.0.0 h1:abc\n")

	bridges1, err := WriteGoModWithVendorBridges(genDir, modRoot, "example.com/a", true, false, "", false, "")
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	if len(bridges1) != 0 {
		t.Fatalf("expected no bridges without vendor, got %d", len(bridges1))
	}
	if _, err := os.Stat(filepath.Join(genDir, gomodSrcFile)); err != nil {
		t.Fatalf("expected %s: %v", gomodSrcFile, err)
	}
	if _, err := os.Stat(filepath.Join(genDir, vendorBridgesCacheFile)); err != nil {
		t.Fatalf("expected %s: %v", vendorBridgesCacheFile, err)
	}
	if err := os.WriteFile(filepath.Join(genDir, "doctest.tidy-done"), []byte("done"), 0644); err != nil {
		t.Fatal(err)
	}
	modPath := filepath.Join(genDir, "go.mod")
	old := time.Now().Add(-3 * time.Hour).UTC().Truncate(time.Second)
	if err := os.Chtimes(modPath, old, old); err != nil {
		t.Fatal(err)
	}
	st1, err := os.Stat(modPath)
	if err != nil {
		t.Fatal(err)
	}

	bridges2, err := WriteGoModWithVendorBridges(genDir, modRoot, "example.com/a", true, false, "", false, "")
	if err != nil {
		t.Fatalf("second: %v", err)
	}
	if len(bridges2) != 0 {
		t.Fatalf("warm bridges: %d", len(bridges2))
	}
	st2, err := os.Stat(modPath)
	if err != nil {
		t.Fatal(err)
	}
	if !st2.ModTime().Equal(st1.ModTime()) {
		t.Fatalf("warm hit must not rewrite go.mod: before=%v after=%v", st1.ModTime(), st2.ModTime())
	}
	if _, err := os.Stat(filepath.Join(genDir, "doctest.tidy-done")); err != nil {
		t.Fatalf("tidy-done must remain on hit: %v", err)
	}
}

func TestWriteGoModModulesTxtInvalidatesSrcCache(t *testing.T) {
	modRoot := t.TempDir()
	genDir := filepath.Join(t.TempDir(), "gen")
	writeTreeFile(t, modRoot, "go.mod", "module example.com/a\n\ngo 1.21\n")
	vendor := filepath.Join(modRoot, "vendor")
	dep := "example.com/dep"
	depDir := filepath.Join(vendor, filepath.FromSlash(dep))
	if err := os.MkdirAll(depDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(depDir, "dep.go"), []byte("package dep\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vendor, "modules.txt"), []byte("# example.com/dep v1.0.0\n## explicit\nexample.com/dep\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := WriteGoModWithVendorBridges(genDir, modRoot, "example.com/a", true, false, "", false, ""); err != nil {
		t.Fatalf("first: %v", err)
	}
	first := readFileString(t, filepath.Join(genDir, "go.mod"))
	if !strings.Contains(first, "example.com/dep") {
		t.Fatalf("expected vendor inject:\n%s", first)
	}
	if err := os.WriteFile(filepath.Join(genDir, "doctest.tidy-done"), []byte("done"), 0644); err != nil {
		t.Fatal(err)
	}

	// modules.txt-only change (go.mod/go.sum untouched).
	if err := os.WriteFile(filepath.Join(vendor, "modules.txt"), []byte("# example.com/dep v1.0.0\n## explicit\nexample.com/dep\n# example.com/other v2.0.0\n## explicit\nexample.com/other\n"), 0644); err != nil {
		t.Fatal(err)
	}
	otherDir := filepath.Join(vendor, "example.com", "other")
	if err := os.MkdirAll(otherDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(otherDir, "other.go"), []byte("package other\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := WriteGoModWithVendorBridges(genDir, modRoot, "example.com/a", true, false, "", false, ""); err != nil {
		t.Fatalf("second after modules.txt: %v", err)
	}
	second := readFileString(t, filepath.Join(genDir, "go.mod"))
	if !strings.Contains(second, "example.com/other") {
		t.Fatalf("expected modules.txt change to regenerate go.mod:\n%s", second)
	}
	if _, err := os.Stat(filepath.Join(genDir, "doctest.tidy-done")); !os.IsNotExist(err) {
		t.Fatalf("tidy-done should drop when gen go.mod rewrote, err=%v", err)
	}
}

func TestWriteGoModVendorBridgeWarmHitReturnsCachedBridges(t *testing.T) {
	modRoot := t.TempDir()
	genDir := filepath.Join(t.TempDir(), "gen")
	writeTreeFile(t, modRoot, "go.mod", "module example.com/a\n\ngo 1.21\n")
	vendor := filepath.Join(modRoot, "vendor")
	nogo := "example.com/nogo"
	nogoDir := filepath.Join(vendor, filepath.FromSlash(nogo))
	if err := os.MkdirAll(nogoDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nogoDir, "nogo.go"), []byte("package nogo\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vendor, "modules.txt"), []byte("# example.com/nogo v0.1.0\n## explicit; go 1.17\nexample.com/nogo\n"), 0644); err != nil {
		t.Fatal(err)
	}

	b1, err := WriteGoModWithVendorBridges(genDir, modRoot, "example.com/a", true, false, "", false, "")
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	if len(b1) != 1 || b1[0].ModulePath != nogo {
		t.Fatalf("first bridges: %#v", b1)
	}
	ph := b1[0].BridgeRoot
	old := time.Now().Add(-2 * time.Hour).UTC().Truncate(time.Second)
	if err := os.Chtimes(ph, old, old); err != nil {
		t.Fatal(err)
	}
	st1, _ := os.Stat(ph)

	b2, err := WriteGoModWithVendorBridges(genDir, modRoot, "example.com/a", true, false, "", false, "")
	if err != nil {
		t.Fatalf("second: %v", err)
	}
	if len(b2) != 1 || b2[0].BridgeRoot != ph {
		t.Fatalf("warm bridges mismatch: %#v vs %#v", b2, b1)
	}
	st2, err := os.Stat(ph)
	if err != nil {
		t.Fatal(err)
	}
	if !st2.ModTime().Equal(st1.ModTime()) {
		t.Fatalf("placeholder must not rewrite on hit")
	}
	if VendorGomodOverlayPath(genDir) == "" {
		t.Fatal("overlay json must remain")
	}
}
