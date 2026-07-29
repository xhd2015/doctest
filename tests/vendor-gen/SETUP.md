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
       [shadow vendor-bridge under genDir when go.mod missing; project vendor read-only]
```

## Preconditions

- Package `github.com/xhd2015/doctest/libdoc/core` is importable.
- Each leaf uses isolated temp `ModRoot` and `GenDir` (`t.TempDir`); never a
  shared mapping-gen cache or the real monorepo vendor tree.
- Classic TDD: `present/*` leaves expect **RED** until WriteGoMod reads
  `vendor/modules.txt` and injects require/replace + placeholders.
- Do **not** read `DOCTEST_SESSION_ID` via `os.Getenv`.

## Steps

1. Root Setup sets default `ModPath` / `HasMod`.
2. Branch Setup prepares either a bare project go.mod or a tiny vendor fixture.
3. Leaf Setup narrows sample module paths / markers for Assert.
4. `Run` calls `core.WriteGoMod` once; Assert inspects gen `go.mod` and vendor
   side effects.

## Context

- Fixture helpers: `prepareFreshDirs`, `seedParentGoMod`, `seedTinyVendor`,
  `writeFile`, `vendorPath`, `hasVendorReplaceLine`, `requireLinePresent`.
- Tiny vendor modules use public-looking paths under `example.com/…` so they
  cannot collide with real network resolution in pure text asserts.
- Placeholder zero version string matches xgo when modules.txt version empty:
  `v0.0.0-00010101000000-000000000000` (only if a leaf seeds empty version).

```go
import (
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

// hasReplaceToVendor reports whether go.mod has a replace of modPath whose
// target contains the project's vendor directory path (filesystem replace).
func hasReplaceToVendor(goMod, modPath, modRoot string) bool {
	vendorPrefix := filepath.Join(modRoot, "vendor")
	for _, line := range strings.Split(goMod, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "replace ") {
			continue
		}
		// replace <modPath> => <path>
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
		left := strings.Join(parts[:arrow], " ")
		// left may be "path" or "path version"
		leftPath := strings.Fields(left)[0]
		if leftPath != modPath {
			continue
		}
		right := parts[arrow+1]
		if right == vendorPrefix || strings.HasPrefix(right, vendorPrefix+string(os.PathSeparator)) ||
			right == filepath.Join(vendorPrefix, filepath.FromSlash(modPath)) ||
			strings.Contains(right, filepath.Join("vendor", filepath.FromSlash(modPath))) {
			return true
		}
		// Also accept slash-normalized contains of /vendor/<modPath>
		slashRight := filepath.ToSlash(right)
		if strings.Contains(slashRight, "/vendor/"+modPath) || strings.HasSuffix(slashRight, "/vendor/"+modPath) {
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
		if strings.Contains(filepath.ToSlash(line), "/vendor/") || strings.Contains(filepath.ToSlash(line), needle) {
			n++
		}
	}
	return n
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
