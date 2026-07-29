## Expected

- `Run` succeeds.
- Gen `go.mod` contains `module testcase` and `replace <project> => <modRoot>`.
- Parent path replace for `example.com/localdep` is present (absolute path).
- **No** replace line targets a path under `modRoot/vendor/` (zero vendor replaces
  from modules.txt — vendor dir does not exist).

## Side Effects

- No `vendor/` directory is created under modRoot by WriteGoMod.

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
	if !strings.Contains(goMod, req.ModRoot) {
		t.Fatalf("expected project replace to point at modRoot %s, got:\n%s", req.ModRoot, goMod)
	}
	// Parent path replace still present (absolute-ized).
	if !strings.Contains(goMod, "replace example.com/localdep =>") {
		t.Fatalf("expected parent path replace for localdep, got:\n%s", goMod)
	}
	if n := countVendorReplaces(goMod, req.ModRoot); n != 0 {
		t.Fatalf("expected 0 vendor replaces when vendor/ absent, got %d:\n%s", n, goMod)
	}
	// Defense in depth: no modules.txt-style sample modules.
	if hasReplaceToVendor(goMod, sampleDepPath, req.ModRoot) {
		t.Fatalf("unexpected vendor replace for %s:\n%s", sampleDepPath, goMod)
	}
	if fileExists(filepath.Join(req.ModRoot, "vendor")) {
		t.Fatal("WriteGoMod must not create vendor/ when it was absent")
	}
}
```
