## Expected

- `Run` succeeds.
- Gen `go.mod` has **exactly one** `replace <SampleModPath> => …` line.
- That replace's RHS is the parent work_around absolute path (or equivalent
  cleaned path under modRoot), **not** `…/vendor/<SampleModPath>`.
- Gen `go.mod` must **not** also contain a vendor-path replace for SampleModPath.
- Other vendor modules not covered by a parent path replace (`example.com/nogo`)
  still get vendor `require` + `replace … => …/vendor/…`.
- Project replace `replace example.com/app => <modRoot>` still present.

## Errors

- Classic TDD: **RED** until WriteGoMod skips vendor inject replace for modules
  that already have a parent filesystem path replace (today dual replace for M).

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

	// --- Parent path replace wins for SampleModPath (overlap with vendor) ---
	repls := replaceRHSes(goMod, req.SampleModPath)
	if len(repls) != 1 {
		t.Fatalf("expected exactly one replace for %s, got %d (RHS=%v):\n%s",
			req.SampleModPath, len(repls), repls, goMod)
	}
	rhs := repls[0]
	// Must not point at vendor tree for this module.
	if hasReplaceToVendor(goMod, req.SampleModPath, req.ModRoot) {
		t.Fatalf("parent path replace must win: no vendor replace for %s, got:\n%s",
			req.SampleModPath, goMod)
	}
	slashRHS := filepath.ToSlash(rhs)
	if strings.Contains(slashRHS, "/vendor/"+req.SampleModPath) ||
		strings.HasSuffix(slashRHS, "/vendor/"+req.SampleModPath) {
		t.Fatalf("replace RHS for %s must not be vendor path, got %q:\n%s",
			req.SampleModPath, rhs, goMod)
	}
	// RHS must be the parent work_around path.
	wantAbs := filepath.Clean(req.ParentLocalAbs)
	if wantAbs == "" {
		t.Fatal("ParentLocalAbs must be set by leaf Setup")
	}
	if filepath.Clean(rhs) != wantAbs && !strings.Contains(filepath.ToSlash(goMod), filepath.ToSlash(wantAbs)) {
		// Accept if the single RHS cleans to wantAbs, or go.mod names wantAbs.
		if filepath.Clean(rhs) != wantAbs {
			t.Fatalf("expected replace %s => parent path %s, got RHS %q:\n%s",
				req.SampleModPath, wantAbs, rhs, goMod)
		}
	}

	// --- Non-overridden vendor modules still get vendor require+replace ---
	if !hasRequire(goMod, req.NoGoModPath, req.NoGoModVersion) {
		t.Fatalf("expected vendor require %s %s for non-overridden module, got:\n%s",
			req.NoGoModPath, req.NoGoModVersion, goMod)
	}
	if !hasReplaceToVendor(goMod, req.NoGoModPath, req.ModRoot) {
		t.Fatalf("expected vendor replace for %s (not in parent replace), got:\n%s",
			req.NoGoModPath, goMod)
	}
}

// replaceRHSes returns the right-hand sides of all replace lines whose left
// module path is modPath (version token on left is ignored if present).
func replaceRHSes(goMod, modPath string) []string {
	var out []string
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
		leftPath := parts[0]
		if leftPath != modPath {
			continue
		}
		out = append(out, parts[arrow+1])
	}
	return out
}
```
