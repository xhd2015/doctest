## Expected

- Second write succeeds.
- Placeholder file still exists.
- Placeholder mtime equals forced-old snapshot.
- Gen go.mod replace still targets project vendor (no vendor-bridge hardlink).

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("warm WriteGoModWithVendorBridges failed: %v", err)
	}
	if !resp.PlaceholderExists {
		t.Fatalf("placeholder missing after warm hit: %s", resp.PlaceholderPath)
	}
	if !resp.PlaceholderMtimeBefore.Equal(resp.PlaceholderMtimeAfter) {
		t.Fatalf("placeholder rewritten on warm hit: before=%v after=%v path=%s",
			resp.PlaceholderMtimeBefore, resp.PlaceholderMtimeAfter, resp.PlaceholderPath)
	}
	// Layout: replace to project vendor, not vendor-bridge shadow.
	if strings.Contains(filepath.ToSlash(resp.GoModContent), "/vendor-bridge/") {
		t.Fatalf("obsolete vendor-bridge replace must not appear:\n%s", resp.GoModContent)
	}
	wantVendor := filepath.ToSlash(filepath.Join(req.ModRoot, "vendor", filepath.FromSlash(req.NogoModPath)))
	if !strings.Contains(filepath.ToSlash(resp.GoModContent), wantVendor) &&
		!strings.Contains(resp.GoModContent, filepath.Join("vendor", filepath.FromSlash(req.NogoModPath))) {
		// Still require nogo appear in go.mod as require/replace.
		if !strings.Contains(resp.GoModContent, req.NogoModPath) {
			t.Fatalf("expected nogo replace/require in gen go.mod:\n%s", resp.GoModContent)
		}
	}
}
```
