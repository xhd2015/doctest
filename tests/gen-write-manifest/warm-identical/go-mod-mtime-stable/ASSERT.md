## Expected

- Second WriteGoMod succeeds.
- `go.mod` mtime equals the pre-second-call forced mtime (no rewrite).
- Unified layout: `doctest.gen-manifest` present, `doctest.gomod-fp` absent.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("second WriteGoMod failed: %v", err)
	}
	requireUnifiedLayout(t, resp)
	if resp.GoModMtimeBefore.IsZero() {
		t.Fatal("missing go.mod mtime snapshot from Setup")
	}
	if !resp.GoModMtimeBefore.Equal(resp.GoModMtimeAfter) {
		t.Fatalf("go.mod mtime changed on identical content: before=%v after=%v",
			resp.GoModMtimeBefore, resp.GoModMtimeAfter)
	}
}
```
