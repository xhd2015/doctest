## Expected

- Suite run succeeds.
- At least two leaf `*_test.go` files found under gen.
- No leaf test file defines `func ExperimentP1RootMarker`.
- No leaf test file contains `type Request` (types live in root package).
- Each leaf test file imports at least one non-stdlib package (the root / ancestor package).

```go
import (
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.RunErr != "" {
		t.Fatalf("ref-mode RunTest failed: %s\nstderr:\n%s", resp.RunErr, resp.Stderr)
	}
	if len(resp.LeafTestFiles) < 2 {
		t.Fatalf("expected ≥2 leaf *_test.go under gen, got %v (all go=%v)",
			basenames(resp.LeafTestFiles), basenames(resp.GoFiles))
	}
	for i, leaf := range resp.LeafTestFiles {
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
