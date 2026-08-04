## Expected

- Second WriteGoMod succeeds.
- `doctest.gomod-src` exists (input fingerprint for warm early-out).
- Unified layout still holds.

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("second WriteGoMod failed: %v", err)
	}
	requireUnifiedLayout(t, resp)
	if !resp.GomodSrcExists {
		t.Fatal("expected doctest.gomod-src after WriteGoMod")
	}
}
```
