## Expected

- `Run` succeeds.
- After WriteGoMod, `vendor/example.com/nogo/go.mod` exists.
- Placeholder contains `module example.com/nogo` and a `go` directive (prefer
  `go 1.17` from modules.txt `## explicit; go 1.17`, or at least any `go X.Y`).
- Existing `vendor/example.com/dep/go.mod` still present (not wiped).

## Errors

- Classic TDD: **RED** — current WriteGoMod does not write vendor placeholders.

## Side Effects

- Placeholder is written under the replace target path (project vendor tree when
  replace points there), so the path is a valid module root for `-mod=mod`.

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
	nogoDir := vendorModuleDir(req, req.NoGoModPath)
	placeholder := filepath.Join(nogoDir, "go.mod")
	if !fileExists(placeholder) {
		t.Fatalf("expected placeholder go.mod at %s after WriteGoMod", placeholder)
	}
	content := readFileOrEmpty(placeholder)
	if !strings.Contains(content, "module "+req.NoGoModPath) {
		t.Fatalf("placeholder must declare module %s, got:\n%s", req.NoGoModPath, content)
	}
	hasGo := false
	for _, line := range strings.Split(content, "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) >= 2 && fields[0] == "go" {
			hasGo = true
			// Prefer modules.txt metadata go 1.17 when implementer preserves it.
			if fields[1] != "1.17" && fields[1] != req.ParentGoVersion && fields[1] != "1.19" {
				// Soft: still accept any go line; log via failure only if missing entirely.
			}
			break
		}
	}
	if !hasGo {
		t.Fatalf("placeholder must include a go directive, got:\n%s", content)
	}
	// Prefer annotated go version from modules.txt when present.
	if !strings.Contains(content, "go 1.17") {
		// Allow fallback to parent go version (implementation detail).
		if req.ParentGoVersion != "" && strings.Contains(content, "go "+req.ParentGoVersion) {
			// ok fallback
		} else {
			t.Fatalf("placeholder go version should be 1.17 (modules.txt) or parent %s, got:\n%s",
				req.ParentGoVersion, content)
		}
	}

	depMod := filepath.Join(vendorModuleDir(req, req.SampleModPath), "go.mod")
	if !fileExists(depMod) {
		t.Fatalf("existing dep go.mod must remain at %s", depMod)
	}
	// Wiring should also reference nogo replace (product completeness).
	if !hasReplaceToVendor(resp.GoModContent, req.NoGoModPath, req.ModRoot) {
		t.Fatalf("gen go.mod should replace %s to vendor path, got:\n%s",
			req.NoGoModPath, resp.GoModContent)
	}
}
```
