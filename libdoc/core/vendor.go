package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// zeroPseudoVersion is used when modules.txt records a module with no version
// (matches xgo / go mod replace-without-require default).
const zeroPseudoVersion = "v0.0.0-00010101000000-000000000000"

// vendorMod is one module entry from vendor/modules.txt.
type vendorMod struct {
	Path           string
	Version        string
	GoVersion      string // from "## …; go X.Y" when present
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

// vendorBridgeForModRoot returns require+replace lines for gen go.mod and
// ensures placeholder go.mod files exist under vendor modules that lack them.
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
func vendorBridgeForModRoot(modRoot, parentGoVersion string, parentPathReplaced, parentModuleReplaced map[string]bool) (requireReplace string, suppressParentModule map[string]bool, err error) {
	suppressParentModule = make(map[string]bool)
	vendorDir := filepath.Join(modRoot, "vendor")
	st, err := os.Stat(vendorDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", suppressParentModule, nil
		}
		return "", suppressParentModule, err
	}
	if !st.IsDir() {
		return "", suppressParentModule, nil
	}
	data, err := os.ReadFile(filepath.Join(vendorDir, "modules.txt"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", suppressParentModule, nil
		}
		return "", suppressParentModule, err
	}
	mods := parseVendorModulesTxt(string(data))
	if len(mods) == 0 {
		return "", suppressParentModule, nil
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

		b.WriteString(fmt.Sprintf("replace %s => %s\n", reqPath, vendorModPath))

		// Placeholder go.mod when missing (xgo-style: module + go).
		goModFile := filepath.Join(vendorModPath, "go.mod")
		if _, serr := os.Stat(goModFile); os.IsNotExist(serr) {
			goVer := m.GoVersion
			if goVer == "" {
				goVer = parentGoVersion
			}
			if goVer == "" {
				goVer = "1.19"
			}
			if err := os.MkdirAll(vendorModPath, 0755); err != nil {
				return "", suppressParentModule, err
			}
			// Blank line between module and go matches common go.mod style;
			// tests only require contains of both directives.
			placeholder := fmt.Sprintf("module %s\n\ngo %s\n", reqPath, goVer)
			if err := os.WriteFile(goModFile, []byte(placeholder), 0644); err != nil {
				return "", suppressParentModule, err
			}
		}
	}
	return b.String(), suppressParentModule, nil
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
