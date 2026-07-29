## Expected

- `Run` succeeds.
- Project `vendor/example.com/nogo/go.mod` must **not** be created (project
  vendor is read-only).
- Gen dir has a shadow under `vendor-bridge/example.com/nogo/` with placeholder
  `go.mod` and hardlinked/copied package sources.
- Gen `go.mod` `replace` for nogo targets the **shadow** path (not project vendor).
- Existing `vendor/example.com/dep/go.mod` still present (not wiped).

## Errors

- Fail if project vendor was mutated, or shadow/placeholder missing.

## Side Effects

- Shadow modules live only under genDir.

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
	projectNogo := filepath.Join(vendorModuleDir(req, req.NoGoModPath), "go.mod")
	if fileExists(projectNogo) {
		t.Fatalf("must not write placeholder into project vendor: %s", projectNogo)
	}

	shadow := filepath.Join(req.GenDir, "vendor-bridge", filepath.FromSlash(req.NoGoModPath))
	ph := filepath.Join(shadow, "go.mod")
	if !fileExists(ph) {
		t.Fatalf("expected shadow go.mod at %s", ph)
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
	// Package source mirrored (hardlink/copy)
	if b, err := os.ReadFile(filepath.Join(shadow, "nogo.go")); err != nil || !strings.Contains(string(b), "package nogo") {
		t.Fatalf("expected package file under shadow: %v %q", err, b)
	}
	// replace targets shadow
	if !strings.Contains(resp.GoModContent, shadow) &&
		!strings.Contains(filepath.ToSlash(resp.GoModContent), "vendor-bridge/"+req.NoGoModPath) {
		t.Fatalf("gen go.mod replace for %s should target vendor-bridge shadow, got:\n%s",
			req.NoGoModPath, resp.GoModContent)
	}

	depMod := filepath.Join(vendorModuleDir(req, req.SampleModPath), "go.mod")
	if !fileExists(depMod) {
		t.Fatalf("existing dep go.mod must remain at %s", depMod)
	}
}
```
