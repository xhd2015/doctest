## Expected

- `Run` succeeds.
- Gen `go.mod` replaces `example.com/dep` to a path under `modRoot/vendor/`.
- The replace target directory contains the distinctive marker string in package
  source (fixture-level proof that gen wiring points at **vendored** content,
  not an empty or cache path).

## Errors

- Classic TDD: **RED** until vendor replace lines are emitted.

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
	if !hasReplaceToVendor(goMod, req.SampleModPath, req.ModRoot) {
		t.Fatalf("expected vendor replace for %s so resolution prefers vendor, got:\n%s",
			req.SampleModPath, goMod)
	}
	// Resolve expected target and confirm marker is still there (and path is the
	// one named by replace when absolute path is used).
	targetDir := vendorModuleDir(req, req.SampleModPath)
	if !strings.Contains(filepath.ToSlash(goMod), filepath.ToSlash(targetDir)) &&
		!strings.Contains(filepath.ToSlash(goMod), "/vendor/"+req.SampleModPath) {
		t.Fatalf("replace must name vendor module dir %s, got:\n%s", targetDir, goMod)
	}
	src := readFileOrEmpty(filepath.Join(targetDir, "dep.go"))
	if !strings.Contains(src, req.DistinctiveMarker) {
		t.Fatalf("replace target source missing marker %q at %s:\n%s",
			req.DistinctiveMarker, targetDir, src)
	}
	// Require line anchors version so tidy cannot freely re-pick cache versions.
	if !hasRequire(goMod, req.SampleModPath, req.SampleModVersion) {
		t.Fatalf("expected require %s %s to pin vendor version, got:\n%s",
			req.SampleModPath, req.SampleModVersion, goMod)
	}
}
```
