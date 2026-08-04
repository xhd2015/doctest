# Scenario

**Feature**: WriteGoMod injects vendor requires/replaces when project has vendor/

```
# project root with go.mod (+ optional vendor/modules.txt)
modRoot
  -> WriteGoMod(genDir, modRoot, modPath, hasMod=true, …)
  -> genDir/go.mod
       module testcase
       go <parent>
       replace <project> => <modRoot>
       [parent path replaces]
       [require+replace each vendor modules.txt entry when vendor/ exists]
       # replace always => project vendor/<mod> (packages stay there)
       # missing go.mod: genDir/vendor-gomod-overlay/<mod>/go.mod + overlay JSON
       # project vendor never written; no package hardlink/copy into gen
```

## Preconditions

- Package `github.com/xhd2015/doctest/libdoc/core` is importable.
- Each leaf uses isolated temp `ModRoot` and `GenDir` (`t.TempDir`); never a
  shared mapping-gen cache or the real monorepo vendor tree.
- Coverage-backfill mode: product already implements xgo-style overlay inject;
  asserts match current correct behavior (GREEN expected).
- Do **not** read `DOCTEST_SESSION_ID` via `os.Getenv`.

## Steps

1. Root Setup sets default `ModPath` / `HasMod`.
2. Branch Setup prepares either a bare project go.mod or a tiny vendor fixture.
3. Leaf Setup narrows sample module paths / markers for Assert.
4. `Run` calls `core.WriteGoMod` once; Assert inspects gen `go.mod`, optional
   overlay artifacts, and project-vendor immutability.

## Context

- Fixture helpers: `prepareFreshDirs`, `seedParentGoMod`, `seedTinyVendor`,
  `writeFile`, `vendorModuleDir`, `hasReplaceToVendor`, `hasRequire`,
  `overlayPlaceholderPath`, `overlayJSONPath`.
- Tiny vendor modules use public-looking paths under `example.com/…` so they
  cannot collide with real network resolution in pure text asserts.
- Placeholder zero version string matches xgo when modules.txt version empty:
  `v0.0.0-00010101000000-000000000000` (only if a leaf seeds empty version).

```go
import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	defaultModPath = "example.com/app"
	// Sample modules for vendor fixtures (stable across present/* leaves).
	sampleDepPath    = "example.com/dep"
	sampleDepVersion = "v1.2.3"
	sampleDepMarker  = "VENDOR_DEP_MARKER_xyz_doctest_p1"
	noGoModPath      = "example.com/nogo"
	noGoModVersion   = "v0.4.0"
	// modules.txt line with empty version is not used by default fixtures.
	vendorGomodOverlayDir  = "vendor-gomod-overlay"
	vendorGomodOverlayJSON = "vendor-gomod-overlay.json"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.ModPath == "" {
		req.ModPath = defaultModPath
	}
	req.HasMod = true
	return nil
}

func prepareFreshDirs(t *testing.T, req *Request) {
	t.Helper()
	req.ModRoot = t.TempDir()
	req.GenDir = filepath.Join(t.TempDir(), "gen")
	if err := os.MkdirAll(req.GenDir, 0755); err != nil {
		t.Fatalf("mkdir gen: %v", err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func seedParentGoMod(t *testing.T, req *Request, body string) {
	t.Helper()
	if req.ModRoot == "" {
		t.Fatal("seedParentGoMod: ModRoot empty")
	}
	if body == "" {
		body = "module " + req.ModPath + "\n\ngo 1.19\n"
	}
	writeFile(t, filepath.Join(req.ModRoot, "go.mod"), body)
	// Capture go version for asserts when body is default-shaped.
	for _, line := range strings.Split(body, "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) >= 2 && fields[0] == "go" {
			req.ParentGoVersion = fields[1]
			break
		}
	}
}

// seedTinyVendor writes vendor/modules.txt and two modules:
//   - example.com/dep@v1.2.3 with go.mod + package containing DistinctiveMarker
//   - example.com/nogo@v0.4.0 with package only (no go.mod) for placeholder leaf
func seedTinyVendor(t *testing.T, req *Request) {
	t.Helper()
	if req.ModRoot == "" {
		t.Fatal("seedTinyVendor: ModRoot empty")
	}
	req.SampleModPath = sampleDepPath
	req.SampleModVersion = sampleDepVersion
	req.NoGoModPath = noGoModPath
	req.NoGoModVersion = noGoModVersion
	if req.DistinctiveMarker == "" {
		req.DistinctiveMarker = sampleDepMarker
	}
	req.VendorRoot = filepath.Join(req.ModRoot, "vendor")

	modulesTxt := strings.Join([]string{
		"# " + sampleDepPath + " " + sampleDepVersion,
		"## explicit; go 1.18",
		sampleDepPath,
		"# " + noGoModPath + " " + noGoModVersion,
		"## explicit; go 1.17",
		noGoModPath,
		"",
	}, "\n")
	writeFile(t, filepath.Join(req.VendorRoot, "modules.txt"), modulesTxt)

	// dep: has go.mod already + distinctive source
	depDir := filepath.Join(req.VendorRoot, filepath.FromSlash(sampleDepPath))
	writeFile(t, filepath.Join(depDir, "go.mod"),
		"module "+sampleDepPath+"\n\ngo 1.18\n")
	writeFile(t, filepath.Join(depDir, "dep.go"),
		"package dep\n\n// "+req.DistinctiveMarker+"\nconst Marker = \""+req.DistinctiveMarker+"\"\n")

	// nogo: package only — no go.mod (placeholder must be created)
	nogoDir := filepath.Join(req.VendorRoot, filepath.FromSlash(noGoModPath))
	writeFile(t, filepath.Join(nogoDir, "nogo.go"),
		"package nogo\n\nconst X = 1\n")
}

func vendorModuleDir(req *Request, modPath string) string {
	return filepath.Join(req.ModRoot, "vendor", filepath.FromSlash(modPath))
}

// overlayPlaceholderPath is genDir/vendor-gomod-overlay/<modPath>/go.mod.
func overlayPlaceholderPath(req *Request, modPath string) string {
	return filepath.Join(req.GenDir, vendorGomodOverlayDir, filepath.FromSlash(modPath), "go.mod")
}

// overlayJSONPath is genDir/vendor-gomod-overlay.json.
func overlayJSONPath(req *Request) string {
	return filepath.Join(req.GenDir, vendorGomodOverlayJSON)
}

// hasReplaceToProjectVendor reports whether go.mod has a filesystem replace of
// modPath targeting project vendor/<modPath> (xgo-style packages stay there).
// Does not accept obsolete vendor-bridge shadow RHS.
func hasReplaceToProjectVendor(goMod, modPath, modRoot string) bool {
	want := filepath.Join(modRoot, "vendor", filepath.FromSlash(modPath))
	wantSlash := filepath.ToSlash(want)
	for _, line := range strings.Split(goMod, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "replace ") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line, "replace "))
		parts := strings.Fields(rest)
		arrow := -1
		for i, p := range parts {
			if p == "=>" {
				arrow = i
				break
			}
		}
		if arrow < 1 || arrow+1 >= len(parts) {
			continue
		}
		leftPath := strings.Fields(strings.Join(parts[:arrow], " "))[0]
		if leftPath != modPath {
			continue
		}
		right := parts[arrow+1]
		slashRight := filepath.ToSlash(right)
		// Reject obsolete vendor-bridge shadows.
		if strings.Contains(slashRight, "/vendor-bridge/") {
			continue
		}
		if right == want || slashRight == wantSlash ||
			strings.HasSuffix(slashRight, "/vendor/"+modPath) {
			return true
		}
	}
	return false
}

// hasReplaceToVendor reports whether go.mod has a filesystem replace of modPath
// targeting project vendor/… (packages). Legacy vendor-bridge RHS is still
// recognized for defensive counting but new leaves prefer hasReplaceToProjectVendor.
func hasReplaceToVendor(goMod, modPath, modRoot string) bool {
	if hasReplaceToProjectVendor(goMod, modPath, modRoot) {
		return true
	}
	for _, line := range strings.Split(goMod, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "replace ") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line, "replace "))
		parts := strings.Fields(rest)
		arrow := -1
		for i, p := range parts {
			if p == "=>" {
				arrow = i
				break
			}
		}
		if arrow < 1 || arrow+1 >= len(parts) {
			continue
		}
		leftPath := strings.Fields(strings.Join(parts[:arrow], " "))[0]
		if leftPath != modPath {
			continue
		}
		slashRight := filepath.ToSlash(parts[arrow+1])
		if strings.Contains(slashRight, "/vendor-bridge/"+modPath) ||
			strings.HasSuffix(slashRight, "/vendor-bridge/"+modPath) {
			return true
		}
	}
	return false
}

func hasRequire(goMod, modPath, version string) bool {
	// require path version  OR  inside require ( ) block as "path version"
	want := modPath + " " + version
	for _, line := range strings.Split(goMod, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == want || strings.HasPrefix(trimmed, "require "+want) {
			return true
		}
	}
	return false
}

func countVendorReplaces(goMod, modRoot string) int {
	n := 0
	needle := filepath.ToSlash(filepath.Join(modRoot, "vendor"))
	for _, line := range strings.Split(goMod, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "replace ") {
			continue
		}
		slash := filepath.ToSlash(line)
		// Project vendor/… (primary). Legacy vendor-bridge still counted if present.
		if strings.Contains(slash, "/vendor/") || strings.Contains(slash, needle) ||
			strings.Contains(slash, "/vendor-bridge/") {
			n++
		}
	}
	return n
}

// parseOverlayReplace reads genDir/vendor-gomod-overlay.json Replace map.
// Returns nil map if file missing or empty.
func parseOverlayReplace(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var payload struct {
		Replace map[string]string `json:"Replace"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return payload.Replace, nil
}

// listNonGoModFiles lists non-go.mod files under dir (non-recursive for leaf files;
// walks for nested package paths). Used to assert no package mirror under overlay.
func listNonGoModFiles(root string) ([]string, error) {
	var out []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Base(path) == "go.mod" {
			return nil
		}
		rel, rerr := filepath.Rel(root, path)
		if rerr != nil {
			out = append(out, path)
		} else {
			out = append(out, rel)
		}
		return nil
	})
	return out, err
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func readFileOrEmpty(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}
```
