package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Source-input early-out for WriteGoModWithVendorBridges (warm prepare).
// Output-side content-hash (doctest.gen-manifest) still applies on miss rebuild.
const (
	gomodSrcFile           = "doctest.gomod-src"
	vendorBridgesCacheFile = "doctest.vendor-bridges.json"
	// Bump when desired gen go.mod shape / inject rules change.
	gomodSrcPolicyVersion = "gomod-src=1"
)

// gomodSrcFingerprint is the canonical multi-line input fingerprint stored in
// doctest.gomod-src. Compare as whole string (exact match).
func gomodSrcFingerprint(modRoot, modPath string, hasMod bool, withAssertReplace bool, assertCacheDir string, withSessionReplace bool, sessionCacheDir string) (string, error) {
	modRootAbs, err := filepath.Abs(modRoot)
	if err != nil {
		modRootAbs = filepath.Clean(modRoot)
	}

	goModH := "n/a"
	goSumH := "n/a"
	if hasMod {
		goModH, err = fileContentHashOrToken(filepath.Join(modRootAbs, "go.mod"))
		if err != nil {
			return "", err
		}
		goSumH, err = fileContentHashOrToken(filepath.Join(modRootAbs, "go.sum"))
		if err != nil {
			return "", err
		}
	}
	modulesH, err := fileContentHashOrToken(filepath.Join(modRootAbs, "vendor", "modules.txt"))
	if err != nil {
		return "", err
	}

	assertKey := "-"
	if withAssertReplace && assertCacheDir != "" && modPath != "github.com/xhd2015/doctest" {
		assertKey = absPathOrRaw(assertCacheDir)
	}
	sessionKey := "-"
	if withSessionReplace && sessionCacheDir != "" && modPath != "github.com/xhd2015/doctest" {
		sessionKey = absPathOrRaw(sessionCacheDir)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "version %s\n", gomodSrcPolicyVersion)
	fmt.Fprintf(&b, "modRoot %s\n", filepath.ToSlash(modRootAbs))
	fmt.Fprintf(&b, "modPath %s\n", modPath)
	fmt.Fprintf(&b, "hasMod %t\n", hasMod)
	fmt.Fprintf(&b, "go.mod %s\n", goModH)
	fmt.Fprintf(&b, "go.sum %s\n", goSumH)
	fmt.Fprintf(&b, "modules.txt %s\n", modulesH)
	fmt.Fprintf(&b, "assert %s\n", assertKey)
	fmt.Fprintf(&b, "session %s\n", sessionKey)
	return b.String(), nil
}

func fileContentHashOrToken(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "absent", nil
		}
		return "", err
	}
	return contentSHA256(data), nil
}

func absPathOrRaw(p string) string {
	if p == "" {
		return "-"
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return filepath.ToSlash(filepath.Clean(p))
	}
	return filepath.ToSlash(abs)
}

// tryGomodSrcHit returns bridges and true when gen root is warm for fp.
// Caller must hold genModMu.
func tryGomodSrcHit(genDir, fp string, hasMod bool, modRoot string) (bridges []VendorBridgeMapping, ok bool) {
	if genDir == "" || fp == "" {
		return nil, false
	}
	stored, err := os.ReadFile(filepath.Join(genDir, gomodSrcFile))
	if err != nil || string(stored) != fp {
		return nil, false
	}
	if _, err := os.Stat(filepath.Join(genDir, "go.mod")); err != nil {
		return nil, false
	}
	// If project has go.sum, gen must too (seeded on last full write).
	if hasMod {
		srcSum := filepath.Join(modRoot, "go.sum")
		if _, err := os.Stat(srcSum); err == nil {
			if _, err := os.Stat(filepath.Join(genDir, "go.sum")); err != nil {
				return nil, false
			}
		}
	}
	bridges, err = loadVendorBridgesCache(genDir)
	if err != nil {
		return nil, false
	}
	for _, b := range bridges {
		if b.BridgeRoot == "" {
			continue
		}
		if _, err := os.Stat(b.BridgeRoot); err != nil {
			return nil, false
		}
	}
	if len(bridges) > 0 && VendorGomodOverlayPath(genDir) == "" {
		return nil, false
	}
	return bridges, true
}

// noteGomodSrcArtifacts marks gomod-src cache + overlay placeholders as desired
// for gen-plan accounting. Shared by warm hit and cold miss so prune/emit stay
// symmetric (placeholders are written with raw WriteFile, not writeRelIfChanged).
func noteGomodSrcArtifacts(genDir string, bridges []VendorBridgeMapping) {
	if genDir == "" {
		return
	}
	NoteDesired(genDir, gomodSrcFile)
	NoteWrite(genDir, gomodSrcFile, false)
	NoteDesired(genDir, vendorBridgesCacheFile)
	NoteWrite(genDir, vendorBridgesCacheFile, false)
	if VendorGomodOverlayPath(genDir) != "" {
		NoteDesired(genDir, VendorGomodOverlayJSON)
		NoteWrite(genDir, VendorGomodOverlayJSON, false)
	}
	for _, b := range bridges {
		if b.BridgeRoot == "" {
			continue
		}
		// Relative under gen when possible for gen-plan desired set.
		if rel, err := filepath.Rel(genDir, b.BridgeRoot); err == nil && rel != "" && !strings.HasPrefix(rel, "..") {
			rel = filepath.ToSlash(rel)
			NoteDesired(genDir, rel)
			NoteWrite(genDir, rel, false)
		}
	}
}

func noteGomodSrcHit(genDir string, bridges []VendorBridgeMapping) {
	noteGenBookkeeping(genDir)
	noteGomodSrcArtifacts(genDir, bridges)
}

type vendorBridgesCachePayload struct {
	Bridges []vendorBridgeCacheEntry `json:"bridges"`
}

type vendorBridgeCacheEntry struct {
	ModulePath         string `json:"modulePath"`
	OriginalVendorRoot string `json:"originalVendorRoot"`
	BridgeRoot         string `json:"bridgeRoot"`
}

func loadVendorBridgesCache(genDir string) ([]VendorBridgeMapping, error) {
	path := filepath.Join(genDir, vendorBridgesCacheFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Missing cache file is only OK when fingerprint also missing (caller
			// checks fp first). Treat as miss.
			return nil, err
		}
		return nil, err
	}
	// Empty bridges is valid (no placeholder modules): file must still exist.
	var payload vendorBridgesCachePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	out := make([]VendorBridgeMapping, 0, len(payload.Bridges))
	for _, e := range payload.Bridges {
		out = append(out, VendorBridgeMapping{
			ModulePath:         e.ModulePath,
			OriginalVendorRoot: e.OriginalVendorRoot,
			BridgeRoot:         e.BridgeRoot,
		})
	}
	return out, nil
}

func saveGomodSrcCache(genDir, fp string, bridges []VendorBridgeMapping) error {
	if genDir == "" {
		return nil
	}
	if err := os.MkdirAll(genDir, 0755); err != nil {
		return err
	}
	srcPath := filepath.Join(genDir, gomodSrcFile)
	if _, err := writeFileIfChanged(srcPath, []byte(fp), 0644); err != nil {
		return err
	}

	entries := make([]vendorBridgeCacheEntry, 0, len(bridges))
	for _, b := range bridges {
		entries = append(entries, vendorBridgeCacheEntry{
			ModulePath:         b.ModulePath,
			OriginalVendorRoot: b.OriginalVendorRoot,
			BridgeRoot:         b.BridgeRoot,
		})
	}
	payload, err := json.Marshal(vendorBridgesCachePayload{Bridges: entries})
	if err != nil {
		return err
	}
	brPath := filepath.Join(genDir, vendorBridgesCacheFile)
	if _, err := writeFileIfChanged(brPath, payload, 0644); err != nil {
		return err
	}
	// Same desired-set notes as warm hit (overlay JSON + placeholders included).
	noteGomodSrcArtifacts(genDir, bridges)
	return nil
}
