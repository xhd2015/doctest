## Expected

```go
import (
	"testing"
	"github.com/xhd2015/doctest/session"
)
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil || resp.ErrMsg != "" { t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg) }
	if resp.OverlayReplace[resp.ProjectSource] != "replacement.go" || resp.OverlayReplace[resp.ActiveBridgeSource] != "replacement.go" { t.Fatalf("overlay=%#v", resp.OverlayReplace) }
	if _, ok := resp.OverlayReplace[resp.ActiveVendorSource]; ok { t.Fatalf("original vendor key remained: %#v", resp.OverlayReplace) }
}
```
