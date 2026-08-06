package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// vendorGomodOverlayDir holds synthetic go.mod files for vendored modules that
// lack one. Packages stay in project vendor/; only go.mod is invented (xgo
// createGoModPlaceholder style).
const vendorGomodOverlayDir = "vendor-gomod-overlay"

// VendorGomodOverlayJSON is the gen-relative Go -overlay file mapping
// project vendor/<mod>/go.mod → placeholder under vendorGomodOverlayDir.
const VendorGomodOverlayJSON = "vendor-gomod-overlay.json"

// zeroPseudoVersion is used when modules.txt records a module with no version
// (matches xgo / go mod replace-without-require default).
const zeroPseudoVersion = "v0.0.0-00010101000000-000000000000"

// vendorMod is one module entry from vendor/modules.txt.
type vendorMod struct {
	Path            string
	Version         string
	GoVersion       string // from "## …; go X.Y" when present
	ReplacementPath string // non-empty when modules.txt records "=> path"
}

// parseVendorModulesTxt parses vendor/modules.txt content (Go modules vendor
// format). Only modules that contribute at least one package line are returned
// (VendorList semantics from cmd/go / xgo goinfo.ParseVendor).
func parseVendorModulesTxt(content string) []vendorMod {
	var (
		list     []vendorMod
		cur      vendorMod
		haveCur  bool
		seenPath = make(map[string]bool) // path → already in list
	)

	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "# ") {
			f := strings.Fields(line)
			if len(f) < 3 {
				haveCur = false
				continue
			}
			cur = vendorMod{}
			// f[0] == "#"
			if isValidSemVer(f[2]) {
				cur.Path = f[1]
				cur.Version = f[2]
				f = f[3:]
			} else if f[2] == "=>" {
				// Wildcard replacement without version: # path => replace
				cur.Path = f[1]
				f = f[2:]
			} else {
				haveCur = false
				continue
			}
			if len(f) >= 2 && f[0] == "=>" {
				if len(f) == 2 {
					cur.ReplacementPath = f[1]
				} else if len(f) >= 3 && isValidSemVer(f[2]) {
					// path version replacement — vendor tree still uses replacement path when non-local
					cur.ReplacementPath = f[1]
				}
			}
			haveCur = cur.Path != ""
			continue
		}

		if !haveCur {
			continue
		}

		if annotations, ok := cutPrefix(line, "## "); ok {
			for _, entry := range strings.Split(annotations, ";") {
				entry = strings.TrimSpace(entry)
				if goVer, ok := cutPrefix(entry, "go "); ok {
					cur.GoVersion = goVer
				}
			}
			continue
		}

		// Package line within current module (single path token).
		if f := strings.Fields(line); len(f) == 1 {
			if !seenPath[cur.Path] {
				list = append(list, cur)
				seenPath[cur.Path] = true
			} else {
				// Update version/metadata if a later package line reaffirms the module
				// (keep first occurrence order; higher version not required for fixtures).
				for i := range list {
					if list[i].Path == cur.Path {
						if cur.Version != "" {
							list[i].Version = cur.Version
						}
						if cur.GoVersion != "" {
							list[i].GoVersion = cur.GoVersion
						}
						if cur.ReplacementPath != "" {
							list[i].ReplacementPath = cur.ReplacementPath
						}
						break
					}
				}
			}
		}
	}
	return list
}

// vendorBridgeForModRoot returns require+replace lines for gen go.mod.
// Modules get replace => project vendor/<mod> (packages stay there). Modules
// missing go.mod get a synthetic go.mod under genDir/vendor-gomod-overlay and
// a -overlay mapping (xgo-style); project vendor/ is never written.
//
// parentGoVersion is e.g. "1.19" (no "go " prefix).
// parentPathReplaced / parentModuleReplaced list module paths that parent
// go.mod already replaces (filesystem path RHS vs module→module RHS).
//
// Parent filesystem path replace always wins (offline-safe already).
// Parent module→module wins unless modules.txt records a non-local private-fork
// form `# A => B` for the same path: then prefer require A + FS replace
// A => vendor/A (packages live under A) and return A in suppressParentModule so
// the caller drops the parent module→module replace (avoids dual-replace and
// bare network require of private B).
//
// Returns empty requireReplace when vendor/ or modules.txt is absent.
// genDir may be empty only when no placeholder go.mod files are needed.
func vendorBridgeForModRoot(modRoot, genDir, parentGoVersion string, parentPathReplaced, parentModuleReplaced map[string]bool) (requireReplace string, suppressParentModule map[string]bool, err error) {
	requireReplace, suppressParentModule, _, err = vendorBridgeForModRootWithMappings(modRoot, genDir, parentGoVersion, parentPathReplaced, parentModuleReplaced)
	return requireReplace, suppressParentModule, err
}

// vendorBridgeForModRootWithMappings is vendorBridgeForModRoot plus metadata for
// modules that needed a synthetic go.mod. VendorBridgeMapping.BridgeRoot is the
// absolute path of the placeholder go.mod file (not a package tree).
func vendorBridgeForModRootWithMappings(modRoot, genDir, parentGoVersion string, parentPathReplaced, parentModuleReplaced map[string]bool) (requireReplace string, suppressParentModule map[string]bool, bridges []VendorBridgeMapping, err error) {
	suppressParentModule = make(map[string]bool)
	vendorDir := filepath.Join(modRoot, "vendor")
	st, err := os.Stat(vendorDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", suppressParentModule, nil, nil
		}
		return "", suppressParentModule, nil, err
	}
	if !st.IsDir() {
		return "", suppressParentModule, nil, nil
	}
	data, err := os.ReadFile(filepath.Join(vendorDir, "modules.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", suppressParentModule, nil, nil
		}
		return "", suppressParentModule, nil, err
	}
	mods := parseVendorModulesTxt(string(data))
	if len(mods) == 0 {
		return "", suppressParentModule, nil, nil
	}

	var b strings.Builder
	for _, m := range mods {
		// Offline-safe private-fork / ordinary vendor graph:
		// modules.txt `# A ver => B ver` still places package lines under A, so
		// gen require/replace left-hand and vendor join path use A (m.Path).
		// Never bare-require non-local B without a filesystem replace for B.
		reqPath := m.Path
		vendorRel := m.Path
		nonLocalFork := m.ReplacementPath != "" && !isLocalFilesystemPath(m.ReplacementPath)

		version := m.Version
		if version == "" {
			version = zeroPseudoVersion
		}
		vendorModPath := filepath.Join(vendorDir, filepath.FromSlash(vendorRel))

		// Always require so the build list can name the module.
		b.WriteString(fmt.Sprintf("require %s %s\n", reqPath, version))

		// Parent filesystem path replace wins: sole replace for this path.
		if parentPathReplaced[reqPath] {
			continue
		}
		// Parent module→module wins unless this is a non-local private-fork
		// entry where vendor FS content under A must win for offline gen.
		if parentModuleReplaced[reqPath] {
			if !nonLocalFork {
				continue
			}
			// Prefer vendor FS replace of A; suppress parent A => B module pin.
			suppressParentModule[reqPath] = true
		}

		// Always replace to project vendor (package sources). Missing go.mod is
		// supplied via -overlay (xgo method), not by hardlinking packages.
		replaceTarget := vendorModPath
		projectGoMod := filepath.Join(vendorModPath, "go.mod")
		if _, serr := os.Stat(projectGoMod); os.IsNotExist(serr) {
			if genDir == "" {
				return "", suppressParentModule, nil, fmt.Errorf("vendor placeholder for %s needs genDir (project vendor is read-only)", reqPath)
			}
			goVer := m.GoVersion
			if goVer == "" {
				goVer = parentGoVersion
			}
			if goVer == "" {
				goVer = "1.19"
			}
			placeholder, perr := writeVendorGoModPlaceholder(genDir, reqPath, goVer)
			if perr != nil {
				return "", suppressParentModule, nil, perr
			}
			bridges = append(bridges, VendorBridgeMapping{
				ModulePath:         reqPath,
				OriginalVendorRoot: vendorModPath,
				BridgeRoot:         placeholder, // absolute placeholder go.mod path
			})
		}

		b.WriteString(fmt.Sprintf("replace %s => %s\n", reqPath, replaceTarget))
	}

	return b.String(), suppressParentModule, bridges, nil
}

// writeVendorGoModPlaceholder writes genDir/vendor-gomod-overlay/<mod>/go.mod
// (content-stable: skips rewrite when unchanged). Returns absolute path.
func writeVendorGoModPlaceholder(genDir, modPath, goVer string) (string, error) {
	dir := filepath.Join(genDir, vendorGomodOverlayDir, filepath.FromSlash(modPath))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, "go.mod")
	content := fmt.Sprintf("module %s\n\ngo %s\n", modPath, goVer)
	if prev, err := os.ReadFile(path); err == nil && string(prev) == content {
		abs, aerr := filepath.Abs(path)
		if aerr != nil {
			return path, nil
		}
		return abs, nil
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path, nil
	}
	return abs, nil
}

// WriteVendorGomodOverlayJSON writes genDir/vendor-gomod-overlay.json for
// `go -overlay=…` (Replace: project vendor/.../go.mod → placeholder).
// When bridges is empty, removes any prior overlay file.
func WriteVendorGomodOverlayJSON(genDir string, bridges []VendorBridgeMapping) error {
	if genDir == "" {
		return nil
	}
	outPath := filepath.Join(genDir, VendorGomodOverlayJSON)
	if len(bridges) == 0 {
		_ = os.Remove(outPath)
		return nil
	}
	replace := make(map[string]string, len(bridges))
	for _, b := range bridges {
		if b.OriginalVendorRoot == "" || b.BridgeRoot == "" {
			continue
		}
		src := filepath.Join(b.OriginalVendorRoot, "go.mod")
		srcAbs, err := filepath.Abs(src)
		if err != nil {
			srcAbs = src
		}
		dstAbs, err := filepath.Abs(b.BridgeRoot)
		if err != nil {
			dstAbs = b.BridgeRoot
		}
		replace[srcAbs] = dstAbs
	}
	if len(replace) == 0 {
		_ = os.Remove(outPath)
		return nil
	}
	payload, err := json.Marshal(struct {
		Replace map[string]string `json:"Replace"`
	}{Replace: replace})
	if err != nil {
		return err
	}
	if prev, err := os.ReadFile(outPath); err == nil && string(prev) == string(payload) {
		return nil
	}
	return os.WriteFile(outPath, payload, 0644)
}

// VendorGomodOverlayPath returns the absolute overlay JSON path when present
// and non-empty; otherwise "".
func VendorGomodOverlayPath(genDir string) string {
	if genDir == "" {
		return ""
	}
	p := filepath.Join(genDir, VendorGomodOverlayJSON)
	st, err := os.Stat(p)
	if err != nil || st.IsDir() || st.Size() == 0 {
		return ""
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return p
	}
	return abs
}

// VendorGomodOverlayGoFlag returns a single "-overlay=PATH" arg, or nil.
func VendorGomodOverlayGoFlag(genDir string) []string {
	p := VendorGomodOverlayPath(genDir)
	if p == "" {
		return nil
	}
	return []string{"-overlay=" + p}
}

// MergeVendorGomodOverlays unions Replace maps from each genRoot's
// vendor-gomod-overlay.json into destDir/vendor-gomod-overlay.json
// (content-stable write). Used by multi-mod hub tidy/test: overlay JSON lives
// under member gen roots, not under __hub.
//
// Returns the absolute path of the merged file when non-empty, or "" when no
// mappings exist (and removes any prior dest overlay file).
func MergeVendorGomodOverlays(destDir string, genRoots []string) (string, error) {
	if destDir == "" {
		return "", nil
	}
	merged := make(map[string]string)
	for _, root := range genRoots {
		if root == "" {
			continue
		}
		p := filepath.Join(root, VendorGomodOverlayJSON)
		data, err := os.ReadFile(p)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", err
		}
		if len(strings.TrimSpace(string(data))) == 0 {
			continue
		}
		var payload struct {
			Replace map[string]string `json:"Replace"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return "", fmt.Errorf("parse %s: %w", p, err)
		}
		for k, v := range payload.Replace {
			if k == "" || v == "" {
				continue
			}
			merged[k] = v
		}
	}
	outPath := filepath.Join(destDir, VendorGomodOverlayJSON)
	if len(merged) == 0 {
		_ = os.Remove(outPath)
		return "", nil
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}
	payload, err := json.Marshal(struct {
		Replace map[string]string `json:"Replace"`
	}{Replace: merged})
	if err != nil {
		return "", err
	}
	if prev, err := os.ReadFile(outPath); err == nil && string(prev) == string(payload) {
		abs, aerr := filepath.Abs(outPath)
		if aerr != nil {
			return outPath, nil
		}
		return abs, nil
	}
	if err := os.WriteFile(outPath, payload, 0644); err != nil {
		return "", err
	}
	abs, err := filepath.Abs(outPath)
	if err != nil {
		return outPath, nil
	}
	return abs, nil
}

// VendorGomodOverlayGoFlagMerged merges member gen-root overlays into destDir
// and returns []string{"-overlay=PATH"} or nil.
func VendorGomodOverlayGoFlagMerged(destDir string, genRoots []string) ([]string, error) {
	p, err := MergeVendorGomodOverlays(destDir, genRoots)
	if err != nil {
		return nil, err
	}
	if p == "" {
		return nil, nil
	}
	return []string{"-overlay=" + p}, nil
}

// MergeReplaceIntoVendorGomodOverlay merges replace keys into
// genDir/vendor-gomod-overlay.json (creates file if needed). Existing keys that
// conflict with a different value return an error; same value is ignored.
// Used for kind A/B internal shims so tidy/test keep using VendorGomodOverlayGoFlag.
func MergeReplaceIntoVendorGomodOverlay(genDir string, replace map[string]string) error {
	if genDir == "" || len(replace) == 0 {
		return nil
	}
	outPath := filepath.Join(genDir, VendorGomodOverlayJSON)
	merged := make(map[string]string)
	if data, err := os.ReadFile(outPath); err == nil && len(strings.TrimSpace(string(data))) > 0 {
		var payload struct {
			Replace map[string]string `json:"Replace"`
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return fmt.Errorf("parse %s: %w", outPath, err)
		}
		for k, v := range payload.Replace {
			merged[k] = v
		}
	}
	for k, v := range replace {
		if k == "" || v == "" {
			continue
		}
		if prev, ok := merged[k]; ok && prev != v {
			return fmt.Errorf("overlay key conflict: %s\n  existing: %s\n  new: %s", k, prev, v)
		}
		merged[k] = v
	}
	if len(merged) == 0 {
		return nil
	}
	if err := os.MkdirAll(genDir, 0755); err != nil {
		return err
	}
	payload, err := json.Marshal(struct {
		Replace map[string]string `json:"Replace"`
	}{Replace: merged})
	if err != nil {
		return err
	}
	if prev, err := os.ReadFile(outPath); err == nil && string(prev) == string(payload) {
		return nil
	}
	return os.WriteFile(outPath, payload, 0644)
}

// AppendGOFLAGSOverlay returns env assignments that set GOFLAGS to include
// -overlay=path (preserves other GOFLAGS tokens from the process environment).
func AppendGOFLAGSOverlay(overlayPath string) []string {
	if overlayPath == "" {
		return nil
	}
	return mergeGOFLAGSOverlay(nil, overlayPath)
}

func isLocalFilesystemPath(p string) bool {
	if p == "" {
		return false
	}
	if p == "." || p == ".." {
		return true
	}
	if strings.HasPrefix(p, "./") || strings.HasPrefix(p, "../") {
		return true
	}
	if filepath.IsAbs(p) {
		return true
	}
	return false
}

func cutPrefix(s, prefix string) (string, bool) {
	if strings.HasPrefix(s, prefix) {
		return s[len(prefix):], true
	}
	return "", false
}

// isValidSemVer is a lightweight check for modules.txt version tokens.
// Accepts v1.2.3, v0.0.0-…, and other go-module version strings starting with 'v'
// plus digits, or pure "vX.Y.Z" shapes. Also accepts versions Go uses without
// the leading v in rare cases — modules.txt always uses go module versions.
func isValidSemVer(v string) bool {
	if v == "" {
		return false
	}
	// Go's modules.txt uses module versions (semver-ish with optional +incompatible).
	// Mirror go's IsValid for practical purposes: start with 'v' and contain a digit.
	if v[0] != 'v' || len(v) < 2 {
		return false
	}
	// Must look like vN… (N digit)
	if v[1] < '0' || v[1] > '9' {
		return false
	}
	// Disallow path-like tokens
	if strings.Contains(v, "/") {
		return false
	}
	return true
}
