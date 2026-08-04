## Expected

- Second write succeeds.
- Gen go.mod mtime equals pre-second forced mtime.
- Cache files still present.

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("warm WriteGoModWithVendorBridges failed: %v", err)
	}
	if resp.GoModMtimeBefore.IsZero() {
		t.Fatal("missing go.mod mtime snapshot from Setup")
	}
	if !resp.GoModMtimeBefore.Equal(resp.GoModMtimeAfter) {
		t.Fatalf("go.mod mtime changed on warm hit: before=%v after=%v",
			resp.GoModMtimeBefore, resp.GoModMtimeAfter)
	}
	requireSrcAndBridges(t, resp)
}
```
