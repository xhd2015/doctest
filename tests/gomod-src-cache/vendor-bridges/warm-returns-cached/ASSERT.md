## Expected

- Second write succeeds.
- Bridge count matches first write.
- First BridgeRoot equals SnapBridgeRoot.
- Module path matches nogo.
- Bridges JSON still on disk; overlay JSON present.

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("warm WriteGoModWithVendorBridges failed: %v", err)
	}
	if resp.BridgeCount != req.SnapBridgeCount {
		t.Fatalf("bridge count mismatch: first=%d warm=%d roots=%v",
			req.SnapBridgeCount, resp.BridgeCount, resp.BridgeRoots)
	}
	if len(resp.BridgeRoots) == 0 || resp.BridgeRoots[0] != req.SnapBridgeRoot {
		t.Fatalf("warm BridgeRoot mismatch: want %q got %v", req.SnapBridgeRoot, resp.BridgeRoots)
	}
	if req.SnapBridgeModulePath != "" {
		if len(resp.BridgeModulePaths) == 0 || resp.BridgeModulePaths[0] != req.SnapBridgeModulePath {
			t.Fatalf("warm module path mismatch: want %q got %v",
				req.SnapBridgeModulePath, resp.BridgeModulePaths)
		}
	}
	requireSrcAndBridges(t, resp)
	if !resp.OverlayJSONExists {
		t.Fatal("overlay json must remain on warm hit with bridges")
	}
}
```
