## Expected

- `Run` succeeds.
- Project `vendor/example.com/nogo/go.mod` must **not** be created (project
  vendor is read-only).
- Gen dir has placeholder only under
  `vendor-gomod-overlay/example.com/nogo/go.mod` (module + go version).
- **No** package sources (e.g. `nogo.go`) under the overlay module dir — no
  hardlink/copy package mirror.
- Gen `go.mod` `replace` for nogo targets **project vendor**, not vendor-bridge
  and not the overlay directory.
- `vendor-gomod-overlay.json` exists with Replace mapping
  `abs(project vendor/.../nogo/go.mod) → abs(overlay placeholder go.mod)`.
- Existing `vendor/example.com/dep/go.mod` still present (not wiped).
- Obsolete `genDir/vendor-bridge/…` package tree must not be required.

## Errors

- Fail if project vendor was mutated, overlay placeholder/JSON missing, package
  files appear under overlay, or replace targets vendor-bridge.

## Side Effects

- Synthetic go.mod files live only under genDir/vendor-gomod-overlay/.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod failed: %v", err)
	}
	projectNogo := filepath.Join(vendorModuleDir(req, req.NoGoModPath), "go.mod")
	if fileExists(projectNogo) {
		t.Fatalf("must not write placeholder into project vendor: %s", projectNogo)
	}

	// Overlay placeholder (not vendor-bridge shadow)
	ph := overlayPlaceholderPath(req, req.NoGoModPath)
	if !fileExists(ph) {
		t.Fatalf("expected overlay placeholder go.mod at %s", ph)
	}
	content := readFileOrEmpty(ph)
	if !strings.Contains(content, "module "+req.NoGoModPath) {
		t.Fatalf("placeholder must declare module %s, got:\n%s", req.NoGoModPath, content)
	}
	if !strings.Contains(content, "go 1.17") &&
		!(req.ParentGoVersion != "" && strings.Contains(content, "go "+req.ParentGoVersion)) {
		t.Fatalf("placeholder go version should be 1.17 or parent %s, got:\n%s",
			req.ParentGoVersion, content)
	}

	// No package mirror under overlay dir
	overlayModDir := filepath.Dir(ph)
	extras, werr := listNonGoModFiles(overlayModDir)
	if werr != nil {
		t.Fatalf("walk overlay mod dir: %v", werr)
	}
	if len(extras) > 0 {
		t.Fatalf("must not hardlink/copy package sources into overlay; found %v under %s",
			extras, overlayModDir)
	}

	// replace targets project vendor (packages stay there)
	if !hasReplaceToProjectVendor(resp.GoModContent, req.NoGoModPath, req.ModRoot) {
		t.Fatalf("gen go.mod replace for %s should target project vendor, got:\n%s",
			req.NoGoModPath, resp.GoModContent)
	}
	if strings.Contains(filepath.ToSlash(resp.GoModContent), "vendor-bridge/") {
		t.Fatalf("obsolete vendor-bridge replace must not appear, got:\n%s", resp.GoModContent)
	}
	if strings.Contains(filepath.ToSlash(resp.GoModContent), "/vendor-gomod-overlay/") {
		t.Fatalf("replace must not point at overlay dir (packages stay in project vendor), got:\n%s",
			resp.GoModContent)
	}

	// Overlay JSON maps phantom project go.mod path → placeholder
	jsonPath := overlayJSONPath(req)
	repl, jerr := parseOverlayReplace(jsonPath)
	if jerr != nil {
		t.Fatalf("parse overlay JSON %s: %v", jsonPath, jerr)
	}
	if len(repl) == 0 {
		t.Fatalf("expected vendor-gomod-overlay.json with Replace mappings at %s", jsonPath)
	}
	wantSrc, _ := filepath.Abs(projectNogo)
	wantDst, _ := filepath.Abs(ph)
	gotDst, ok := repl[wantSrc]
	if !ok {
		// tolerate non-abs keys if product wrote relative (should be abs)
		for k, v := range repl {
			if filepath.Clean(k) == filepath.Clean(wantSrc) ||
				strings.HasSuffix(filepath.ToSlash(k), "/vendor/"+req.NoGoModPath+"/go.mod") {
				gotDst, ok = v, true
				break
			}
		}
	}
	if !ok {
		t.Fatalf("overlay Replace missing key for %s; map=%v", wantSrc, repl)
	}
	if filepath.Clean(gotDst) != filepath.Clean(wantDst) &&
		!strings.HasSuffix(filepath.ToSlash(gotDst), "/vendor-gomod-overlay/"+req.NoGoModPath+"/go.mod") {
		t.Fatalf("overlay Replace[%s] = %q, want placeholder %s", wantSrc, gotDst, wantDst)
	}

	// Obsolete vendor-bridge tree is not required (anti-regression: must not be the design)
	bridgeShadow := filepath.Join(req.GenDir, "vendor-bridge", filepath.FromSlash(req.NoGoModPath))
	if fileExists(filepath.Join(bridgeShadow, "go.mod")) && !fileExists(ph) {
		t.Fatalf("legacy vendor-bridge without overlay is obsolete")
	}

	depMod := filepath.Join(vendorModuleDir(req, req.SampleModPath), "go.mod")
	if !fileExists(depMod) {
		t.Fatalf("existing dep go.mod must remain at %s", depMod)
	}
}
```
