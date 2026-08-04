## Expected

- `Run` succeeds.
- Gen `go.mod` has `require` + `replace … => project vendor` for `example.com/dep`.
- No directory `genDir/vendor-gomod-overlay/example.com/dep/`.
- No `genDir/vendor-gomod-overlay.json` (or empty/absent — no placeholder bridges).
- Project `vendor/example.com/dep/go.mod` content unchanged / still present.
- No `vendor-bridge` replace targets.

## Errors

- Fail if overlay placeholder or overlay JSON is written when every modules.txt
  module already has go.mod under project vendor.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod failed: %v", err)
	}
	goMod := resp.GoModContent
	if !hasRequire(goMod, req.SampleModPath, req.SampleModVersion) {
		t.Fatalf("expected require %s %s, got:\n%s",
			req.SampleModPath, req.SampleModVersion, goMod)
	}
	if !hasReplaceToProjectVendor(goMod, req.SampleModPath, req.ModRoot) {
		t.Fatalf("expected replace %s => project vendor, got:\n%s",
			req.SampleModPath, goMod)
	}
	if strings.Contains(filepath.ToSlash(goMod), "vendor-bridge/") {
		t.Fatalf("obsolete vendor-bridge replace must not appear, got:\n%s", goMod)
	}

	// No overlay placeholder for modules that already have go.mod
	overlayDep := filepath.Join(req.GenDir, vendorGomodOverlayDir, filepath.FromSlash(req.SampleModPath))
	if fileExists(overlayDep) || fileExists(filepath.Join(overlayDep, "go.mod")) {
		t.Fatalf("must not create overlay placeholder when project go.mod exists: %s", overlayDep)
	}

	// No overlay JSON when no placeholders were needed
	jsonPath := overlayJSONPath(req)
	if fileExists(jsonPath) {
		st, _ := os.Stat(jsonPath)
		if st != nil && st.Size() > 0 {
			repl, jerr := parseOverlayReplace(jsonPath)
			if jerr != nil {
				t.Fatalf("overlay JSON present but unreadable: %v", jerr)
			}
			if len(repl) > 0 {
				t.Fatalf("expected no overlay Replace mappings when all modules have go.mod, got %v", repl)
			}
		}
	}

	depMod := filepath.Join(vendorModuleDir(req, req.SampleModPath), "go.mod")
	if !fileExists(depMod) {
		t.Fatalf("existing dep go.mod must remain at %s", depMod)
	}
}
```
