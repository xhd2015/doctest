## Expected

- Only the package overlay from the hook remains (no phantom go.mod merge without metadata).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil || resp.ErrMsg != "" {
		t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg)
	}
	if resp.OverlayReplace[resp.ActiveVendorSource] != "replacement.go" || len(resp.OverlayReplace) != 1 {
		t.Fatalf("stale placeholder must not merge go.mod: %#v", resp.OverlayReplace)
	}
	if _, ok := resp.OverlayReplace[resp.ActiveVendorGoMod]; ok {
		t.Fatalf("unexpected go.mod key without active metadata: %#v", resp.OverlayReplace)
	}
}
```
