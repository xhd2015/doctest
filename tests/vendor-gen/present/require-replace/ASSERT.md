## Expected

- `Run` succeeds.
- Gen `go.mod` includes `require example.com/dep v1.2.3` (or equivalent require
  block form) and `replace example.com/dep =>` targeting
  `<modRoot>/vendor/example.com/dep`.
- Same require for `example.com/nogo` at `v0.4.0`; its replace targets gen
  `vendor-bridge/…` (shadow; no project `go.mod`).
- Project replace `replace example.com/app => <modRoot>` still present.

## Errors

- Classic TDD: **RED** until WriteGoMod parses modules.txt and emits require +
  replace lines (today only project/parent replaces are written).

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
	goMod := resp.GoModContent
	if !strings.Contains(goMod, "module testcase") {
		t.Fatalf("expected module testcase, got:\n%s", goMod)
	}
	if !strings.Contains(goMod, "replace "+req.ModPath+" =>") {
		t.Fatalf("expected project replace, got:\n%s", goMod)
	}

	// Sample dep: require + replace -> vendor
	if !hasRequire(goMod, req.SampleModPath, req.SampleModVersion) {
		t.Fatalf("expected require %s %s, got:\n%s",
			req.SampleModPath, req.SampleModVersion, goMod)
	}
	if !hasReplaceToVendor(goMod, req.SampleModPath, req.ModRoot) {
		t.Fatalf("expected replace %s => %s/vendor/..., got:\n%s",
			req.SampleModPath, req.ModRoot, goMod)
	}

	// Second module from modules.txt (no project go.mod → vendor-bridge shadow)
	if !hasRequire(goMod, req.NoGoModPath, req.NoGoModVersion) {
		t.Fatalf("expected require %s %s, got:\n%s",
			req.NoGoModPath, req.NoGoModVersion, goMod)
	}
	if !strings.Contains(filepath.ToSlash(goMod), "vendor-bridge/"+req.NoGoModPath) {
		t.Fatalf("expected replace %s => vendor-bridge shadow, got:\n%s",
			req.NoGoModPath, goMod)
	}

	// At least two vendor replaces from our fixture.
	if n := countVendorReplaces(goMod, req.ModRoot); n < 2 {
		t.Fatalf("expected >=2 vendor replaces, got %d:\n%s", n, goMod)
	}

	// Replace target path should resolve to the seeded vendor module dir.
	depVendor := vendorModuleDir(req, req.SampleModPath)
	if !strings.Contains(filepath.ToSlash(goMod), filepath.ToSlash(depVendor)) &&
		!strings.Contains(filepath.ToSlash(goMod), "/vendor/"+req.SampleModPath) {
		t.Fatalf("replace target should reference %s, got:\n%s", depVendor, goMod)
	}
}
```
