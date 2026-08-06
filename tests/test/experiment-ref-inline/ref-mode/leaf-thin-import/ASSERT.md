## Expected

- Suite run succeeds.
- At least two leaf package `.go` files found under gen (`a/` / `b/`).
- No leaf package file defines `func ExperimentP1RootMarker`.
- No leaf package file contains `type Request` (types live in root package).
- Each leaf package file imports at least one non-stdlib package (root / ancestor).

```go
import (
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.RunErr != "" {
		t.Fatalf("hierarchical RunTest failed: %s\nstderr:\n%s", resp.RunErr, resp.Stderr)
	}
	if len(resp.LeafGoFiles) < 2 {
		t.Fatalf("expected ≥2 leaf .go under gen a/b, got %v (all go=%v)",
			basenames(resp.LeafGoFiles), basenames(resp.GoFiles))
	}
	for i, leaf := range resp.LeafGoFiles {
		if i < len(resp.LeafHasMarkerDef) && resp.LeafHasMarkerDef[i] {
			t.Fatalf("thin leaf %s must not define %s", leaf, experimentP1MarkerFunc)
		}
		if i < len(resp.LeafHasTypeReq) && resp.LeafHasTypeReq[i] {
			t.Fatalf("thin leaf %s must not inline type Request (want import of root package)", leaf)
		}
		var imports []string
		if i < len(resp.LeafImportLines) {
			imports = resp.LeafImportLines[i]
		}
		if !leafHasNonStdImport(imports) {
			t.Fatalf("thin leaf %s must import root/ancestor package; imports=%v",
				filepath.ToSlash(leaf), imports)
		}
	}
}

func basenames(paths []string) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		out[i] = filepath.Base(p)
	}
	return out
}
```
