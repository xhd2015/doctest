## Expected

- Second WriteGoMod succeeds.
- Generated go.mod no longer contains `replace localdep`.
- `doctest.gen-manifest` exists and lists `go.mod`.
- Manifest entry for `go.mod` differs from the pre-change snapshot (hash updated).
- No `doctest.gomod-fp`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod after source change failed: %v", err)
	}
	requireUnifiedLayout(t, resp)
	if strings.Contains(resp.GoModContent, "replace localdep") {
		t.Fatalf("expected regenerated go.mod to drop replace localdep, got:\n%s", resp.GoModContent)
	}
	if !strings.Contains(snapGoModContentBefore, "replace localdep") {
		t.Fatalf("precondition failed: first go.mod should have contained replace localdep:\n%s", snapGoModContentBefore)
	}
	if resp.ManifestEntryAfter == "" {
		t.Fatalf("manifest missing go.mod entry after change:\n%s", resp.ManifestContent)
	}
	if snapManifestEntryBefore != "" && resp.ManifestEntryAfter == snapManifestEntryBefore {
		t.Fatalf("expected manifest go.mod entry to change when content changed:\n%s", resp.ManifestEntryAfter)
	}
	if snapManifestEntryBefore == "" {
		// First write under RED may lack manifest; after GREEN first write always sets it.
		// Still require a non-empty entry after the change path.
		t.Log("note: no pre-change manifest entry (first write lacked unified manifest)")
	}
}
```
