## Expected

- Cold WriteGoModWithVendorBridges succeeds.
- `doctest.gomod-src` exists.
- `doctest.vendor-bridges.json` exists (empty bridges array OK with no vendor).
- Bridge return count is 0 when no vendor placeholders.
- Gen `go.mod` exists (module testcase).

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("cold WriteGoModWithVendorBridges failed: %v", err)
	}
	requireSrcAndBridges(t, resp)
	if resp.BridgeCount != 0 {
		t.Fatalf("expected empty bridges without vendor, got %d roots=%v", resp.BridgeCount, resp.BridgeRoots)
	}
	if strings.TrimSpace(resp.BridgesJSONContent) == "" {
		t.Fatal("vendor-bridges.json must not be empty")
	}
	if !strings.Contains(resp.BridgesJSONContent, "bridges") {
		t.Fatalf("vendor-bridges.json should include bridges key, got: %s", resp.BridgesJSONContent)
	}
	if !fileExists(filepath.Join(req.GenDir, "go.mod")) {
		t.Fatal("expected gen go.mod after cold write")
	}
	if !strings.Contains(resp.GoModContent, "module testcase") {
		t.Fatalf("expected gen module testcase, got:\n%s", resp.GoModContent)
	}
}
```
