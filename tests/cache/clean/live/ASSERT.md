## Expected

- Exit code 0.
- Stdout contains `Removed` and the absolute DoctestRoot path.
- Human size present on the remove line.
- **Side effect**: DoctestRoot no longer exists (entire tree removed).
- CacheHome parent may still exist (only `doctest` component is wiped).

## Side Effects

- `$CacheHome/doctest` tree is gone after a successful clean.

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
	mustContain(t, out, "Removed", "Removed prefix")
	mustContain(t, out, req.DoctestRoot, "doctest root path")
	if !hasHumanSizeUnit(out) {
		t.Fatalf("expected human size on Removed line:\n%s", out)
	}
	if pathExists(req.DoctestRoot) {
		t.Fatalf("live clean must remove DoctestRoot: %s\nstdout:\n%s", req.DoctestRoot, out)
	}
	// Parent CacheHome should still exist (we only wipe the doctest component).
	if req.CacheHome != "" && !pathExists(req.CacheHome) {
		t.Fatalf("clean should not remove CacheHome itself: %s", req.CacheHome)
	}
}
```
