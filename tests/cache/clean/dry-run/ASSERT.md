## Expected

- Exit code 0.
- Stdout contains `[dry-run] would remove:` and the absolute DoctestRoot path.
- Line includes a human size in parentheses (or adjacent).
- **Side effect**: DoctestRoot still exists after the run; seeded content remains.

## Side Effects

- No deletes under the injectable CacheHome/doctest tree.

## Exit Code

- 0

```go
import (
	"path/filepath"
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
	mustContain(t, out, req.DoctestRoot, "doctest root path")
	if !hasHumanSizeUnit(out) {
		t.Fatalf("expected human size on dry-run line:\n%s", out)
	}
	if !pathExists(req.DoctestRoot) {
		t.Fatalf("dry-run must not delete DoctestRoot: %s", req.DoctestRoot)
	}
	bucket := filepath.Join(req.DoctestRoot, "leaf-cache")
	if !pathExists(bucket) {
		t.Fatalf("dry-run must not delete seeded bucket leaf-cache under %s", req.DoctestRoot)
	}
}
```
