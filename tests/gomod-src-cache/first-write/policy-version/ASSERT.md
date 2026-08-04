## Expected

- Write succeeds.
- `doctest.gomod-src` content starts with `version gomod-src=1`.
- Fingerprint mentions go.mod / modules.txt keys (format smoke).
- `doctest.gomod-fp` is absent.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("WriteGoModWithVendorBridges failed: %v", err)
	}
	if !resp.GomodSrcExists {
		t.Fatal("expected doctest.gomod-src")
	}
	if !fingerprintHasPolicyVersion(resp.GomodSrcContent) {
		t.Fatalf("expected fingerprint to start with version gomod-src=1, got:\n%s", resp.GomodSrcContent)
	}
	if !strings.Contains(resp.GomodSrcContent, "go.mod ") {
		t.Fatalf("fingerprint should include go.mod hash line:\n%s", resp.GomodSrcContent)
	}
	if !strings.Contains(resp.GomodSrcContent, "modules.txt ") {
		t.Fatalf("fingerprint should include modules.txt hash line:\n%s", resp.GomodSrcContent)
	}
	if resp.GomodFpExists {
		t.Fatal("legacy doctest.gomod-fp must not exist")
	}
}
```
