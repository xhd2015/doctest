package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVendorBridgeShadowNotProjectVendor(t *testing.T) {
	modRoot := t.TempDir()
	genDir := t.TempDir()
	vendor := filepath.Join(modRoot, "vendor")
	nogoPath := "example.com/nogo"
	nogoDir := filepath.Join(vendor, filepath.FromSlash(nogoPath))
	if err := os.MkdirAll(nogoDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nogoDir, "nogo.go"), []byte("package nogo\nconst X=1\n"), 0644); err != nil {
		t.Fatal(err)
	}
	modules := "# example.com/nogo v0.1.0\n## explicit; go 1.17\nexample.com/nogo\n"
	if err := os.WriteFile(filepath.Join(vendor, "modules.txt"), []byte(modules), 0644); err != nil {
		t.Fatal(err)
	}

	extra, _, err := vendorBridgeForModRoot(modRoot, genDir, "1.19", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(extra, "require example.com/nogo") {
		t.Fatalf("require missing: %s", extra)
	}
	// Project vendor must not gain go.mod
	if _, err := os.Stat(filepath.Join(nogoDir, "go.mod")); !os.IsNotExist(err) {
		t.Fatalf("project vendor go.mod must not exist, err=%v", err)
	}
	// Shadow under gen
	shadow := filepath.Join(genDir, vendorBridgeDir, filepath.FromSlash(nogoPath))
	ph, err := os.ReadFile(filepath.Join(shadow, "go.mod"))
	if err != nil {
		t.Fatalf("shadow go.mod: %v", err)
	}
	if !strings.Contains(string(ph), "module example.com/nogo") || !strings.Contains(string(ph), "go 1.17") {
		t.Fatalf("placeholder: %s", ph)
	}
	// Package file mirrored (hardlink or copy — not a bare missing path)
	mirrored := filepath.Join(shadow, "nogo.go")
	data, err := os.ReadFile(mirrored)
	if err != nil {
		t.Fatalf("mirrored nogo.go: %v", err)
	}
	if !strings.Contains(string(data), "package nogo") {
		t.Fatalf("mirrored content: %s", data)
	}
	// replace points at shadow, not project vendor
	if !strings.Contains(extra, "replace example.com/nogo => "+shadow) &&
		!strings.Contains(filepath.ToSlash(extra), "vendor-bridge/example.com/nogo") {
		t.Fatalf("replace should target shadow, got:\n%s", extra)
	}
}

func TestVendorBridgeUsesProjectWhenGoModExists(t *testing.T) {
	modRoot := t.TempDir()
	genDir := t.TempDir()
	vendor := filepath.Join(modRoot, "vendor")
	dep := "example.com/dep"
	depDir := filepath.Join(vendor, filepath.FromSlash(dep))
	if err := os.MkdirAll(depDir, 0755); err != nil {
		t.Fatal(err)
	}
	existing := "module example.com/dep\n\ngo 1.18\n"
	if err := os.WriteFile(filepath.Join(depDir, "go.mod"), []byte(existing), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(depDir, "dep.go"), []byte("package dep\n"), 0644); err != nil {
		t.Fatal(err)
	}
	modules := "# example.com/dep v1.0.0\n## explicit\nexample.com/dep\n"
	if err := os.WriteFile(filepath.Join(vendor, "modules.txt"), []byte(modules), 0644); err != nil {
		t.Fatal(err)
	}

	extra, _, err := vendorBridgeForModRoot(modRoot, genDir, "1.19", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(filepath.Join(depDir, "go.mod"))
	if string(got) != existing {
		t.Fatalf("existing go.mod mutated: %q", got)
	}
	// No shadow for modules that already have go.mod
	if _, err := os.Stat(filepath.Join(genDir, vendorBridgeDir, filepath.FromSlash(dep))); !os.IsNotExist(err) {
		t.Fatalf("should not create shadow when project go.mod exists")
	}
	if !strings.Contains(extra, depDir) && !strings.Contains(filepath.ToSlash(extra), "/vendor/example.com/dep") {
		t.Fatalf("replace should point at project vendor, got:\n%s", extra)
	}
}
