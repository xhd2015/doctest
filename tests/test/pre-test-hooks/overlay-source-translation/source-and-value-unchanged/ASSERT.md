## Expected

- Project and package overlay keys keep their replacement values.
- Phantom go.mod pair is merged without rewriting package keys.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil || resp.ErrMsg != "" {
		t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg)
	}
	if resp.OverlayReplace[resp.ProjectSource] != "replacement.go" ||
		resp.OverlayReplace[resp.ActiveVendorSource] != "replacement.go" {
		t.Fatalf("project/package keys or values changed: %#v", resp.OverlayReplace)
	}
	if got := resp.OverlayReplace[resp.ActiveVendorGoMod]; got != resp.ActiveBridgeSource {
		t.Fatalf("phantom go.mod pair: got=%q want=%q overlay=%#v", got, resp.ActiveBridgeSource, resp.OverlayReplace)
	}
}
```
