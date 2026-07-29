## Expected

- `Run` succeeds.
- Gen `go.mod` contains `go 1.20` (parent directive preserved).
- Project replace to modRoot present.
- Parent path replace `example.com/localdep` present (absolute path under modRoot).
- Vendor require + replace for `example.com/dep` present alongside the above.

## Errors

- Classic TDD: **RED** because vendor require/replace is missing today; parent
  go + path replace halves may already pass in isolation.

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

	// Parent go directive
	if !strings.Contains(goMod, "go 1.20") {
		t.Fatalf("expected parent go 1.20 in gen go.mod, got:\n%s", goMod)
	}

	// Project replace
	if !strings.Contains(goMod, "replace "+req.ModPath+" =>") || !strings.Contains(goMod, req.ModRoot) {
		t.Fatalf("expected project replace => modRoot, got:\n%s", goMod)
	}

	// Parent path replace (absolute-ized by existing WriteGoMod)
	if !strings.Contains(goMod, "replace example.com/localdep =>") {
		t.Fatalf("expected parent path replace for localdep, got:\n%s", goMod)
	}
	if req.ParentLocalAbs != "" && !strings.Contains(goMod, req.ParentLocalAbs) {
		// Accept cleaned absolute path variants
		clean := filepath.Clean(req.ParentLocalAbs)
		if !strings.Contains(goMod, clean) {
			t.Fatalf("expected localdep replace to absolute %s, got:\n%s", clean, goMod)
		}
	}

	// Vendor wiring coexists
	if !hasRequire(goMod, req.SampleModPath, req.SampleModVersion) {
		t.Fatalf("expected vendor require %s %s alongside parent replaces, got:\n%s",
			req.SampleModPath, req.SampleModVersion, goMod)
	}
	if !hasReplaceToVendor(goMod, req.SampleModPath, req.ModRoot) {
		t.Fatalf("expected vendor replace for %s alongside parent replaces, got:\n%s",
			req.SampleModPath, goMod)
	}
}
```
