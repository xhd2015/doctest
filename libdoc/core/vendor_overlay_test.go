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

	extra, _, bridges, err := vendorBridgeForModRootWithMappings(modRoot, genDir, "1.19", nil, nil)
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
	// xgo-style: placeholder only under vendor-gomod-overlay (no package hardlinks)
	phPath := filepath.Join(genDir, vendorGomodOverlayDir, filepath.FromSlash(nogoPath), "go.mod")
	ph, err := os.ReadFile(phPath)
	if err != nil {
		t.Fatalf("placeholder go.mod: %v", err)
	}
	if !strings.Contains(string(ph), "module example.com/nogo") || !strings.Contains(string(ph), "go 1.17") {
		t.Fatalf("placeholder: %s", ph)
	}
	// Package must NOT be mirrored into gen
	if _, err := os.Stat(filepath.Join(genDir, vendorGomodOverlayDir, filepath.FromSlash(nogoPath), "nogo.go")); !os.IsNotExist(err) {
		t.Fatalf("must not hardlink/copy package sources into overlay dir")
	}
	// replace points at project vendor (packages live there)
	if !strings.Contains(extra, "replace example.com/nogo => "+nogoDir) &&
		!strings.Contains(filepath.ToSlash(extra), "/vendor/example.com/nogo") {
		t.Fatalf("replace should target project vendor, got:\n%s", extra)
	}
	if len(bridges) != 1 || bridges[0].BridgeRoot == "" {
		t.Fatalf("expected one bridge mapping to placeholder, got %#v", bridges)
	}
	if err := WriteVendorGomodOverlayJSON(genDir, bridges); err != nil {
		t.Fatal(err)
	}
	if VendorGomodOverlayPath(genDir) == "" {
		t.Fatal("expected vendor-gomod-overlay.json")
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
	// No placeholder for modules that already have go.mod
	if _, err := os.Stat(filepath.Join(genDir, vendorGomodOverlayDir, filepath.FromSlash(dep))); !os.IsNotExist(err) {
		t.Fatalf("should not create placeholder when project go.mod exists")
	}
	if !strings.Contains(extra, depDir) && !strings.Contains(filepath.ToSlash(extra), "/vendor/example.com/dep") {
		t.Fatalf("replace should point at project vendor, got:\n%s", extra)
	}
}
