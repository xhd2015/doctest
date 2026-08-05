## Expected

- Exit code 0.
- Stdout dry-run lines include **both**:
  - main `DoctestRoot` (`$CacheHome/doctest`)
  - outside `LeafCache` path
- Each line uses `[dry-run] would remove:` style.
- **Side effects**: both paths still exist (dry-run only).

## Side Effects

- No deletes of DoctestRoot or LeafCache.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	requireOK(t, resp, err)
	out := resp.Stdout
	if strings.Contains(combinedOut(resp), "unknown command: cache") {
		t.Fatalf("cache not implemented yet:\n%s", combinedOut(resp))
	}
	mustContain(t, out, "[dry-run]", "dry-run prefix")
	mustContain(t, out, "would remove", "would remove")
	mustContain(t, out, req.DoctestRoot, "main doctest root")
	if req.LeafCache == "" {
		t.Fatal("req.LeafCache must be set by Setup")
	}
	mustContain(t, out, req.LeafCache, "outside leaf-cache override")
	if !pathExists(req.DoctestRoot) {
		t.Fatalf("dry-run must not delete DoctestRoot: %s", req.DoctestRoot)
	}
	if !pathExists(req.LeafCache) {
		t.Fatalf("dry-run must not delete LeafCache: %s", req.LeafCache)
	}
}
```
