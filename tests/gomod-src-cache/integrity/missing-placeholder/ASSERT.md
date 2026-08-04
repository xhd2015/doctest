## Expected

- Second write succeeds after placeholder deletion.
- Placeholder path exists again (restored on rebuild).
- Returned bridges non-empty; module path is nogo.
- Overlay JSON present when bridges non-empty.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoMod after deleting placeholder failed: %v", err)
	}
	if !resp.PlaceholderExists {
		t.Fatalf("expected placeholder restored at %s", resp.PlaceholderPath)
	}
	if req.SnapPlaceholderPath != "" && !fileExists(req.SnapPlaceholderPath) {
		// Rebuild may rewrite same path; if path changed, still require some placeholder.
		if resp.BridgeCount == 0 {
			t.Fatalf("placeholder %s not restored and no bridges returned", req.SnapPlaceholderPath)
		}
	}
	if resp.BridgeCount == 0 {
		t.Fatal("expected bridges after nogo rebuild")
	}
	foundNogo := false
	for _, p := range resp.BridgeModulePaths {
		if p == req.NogoModPath || strings.Contains(p, "nogo") {
			foundNogo = true
			break
		}
	}
	if !foundNogo {
		t.Fatalf("expected nogo in bridge modules, got %v", resp.BridgeModulePaths)
	}
	if !resp.OverlayJSONExists {
		t.Fatal("expected vendor-gomod-overlay.json when bridges non-empty")
	}
	requireSrcAndBridges(t, resp)
}
```
