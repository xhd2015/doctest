## Expected

```go
import (
	"testing"
	"github.com/xhd2015/doctest/session"
)
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil || resp.ErrMsg != "" { t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg) }
	if got := resp.OverlayReplace[resp.ActiveBridgeSource]; got != "replacement.go" { t.Fatalf("bridge replacement=%q overlay=%#v", got, resp.OverlayReplace) }
	if _, ok := resp.OverlayReplace[resp.ActiveVendorSource]; ok { t.Fatalf("original vendor key remained: %#v", resp.OverlayReplace) }
}
```
