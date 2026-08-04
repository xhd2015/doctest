## Expected

- Second write succeeds.
- `doctest.gomod-src` content equals pre-second snapshot.
- Fingerprint still has policy version line.

```go
import (
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("warm WriteGoModWithVendorBridges failed: %v", err)
	}
	if !resp.GomodSrcExists {
		t.Fatal("expected doctest.gomod-src after warm hit")
	}
	if resp.GomodSrcContent != req.SnapGomodSrcBefore {
		t.Fatalf("gomod-src content changed on warm hit:\nbefore=%q\nafter=%q",
			req.SnapGomodSrcBefore, resp.GomodSrcContent)
	}
	if !fingerprintHasPolicyVersion(resp.GomodSrcContent) {
		t.Fatalf("policy version lost:\n%s", resp.GomodSrcContent)
	}
}
```
