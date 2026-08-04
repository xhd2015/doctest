## Expected

- Package overlay key remains project `vendor/.../pkg/active.go`.
- Phantom go.mod pair is merged: vendor `.../go.mod` → placeholder go.mod path.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil || resp.ErrMsg != "" {
		t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg)
	}
	if got := resp.OverlayReplace[resp.ActiveVendorSource]; got != "replacement.go" {
		t.Fatalf("package key should stay on project vendor: got=%q overlay=%#v", got, resp.OverlayReplace)
	}
	if got := resp.OverlayReplace[resp.ActiveVendorGoMod]; got != resp.ActiveBridgeSource {
		t.Fatalf("phantom go.mod pair: got=%q want=%q overlay=%#v", got, resp.ActiveBridgeSource, resp.OverlayReplace)
	}
}
```
