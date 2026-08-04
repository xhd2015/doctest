## Expected

```go
import (
	"testing"
	"github.com/xhd2015/doctest/session"
)
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil || resp.ErrMsg != "" { t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg) }
	if resp.OverlayReplace[resp.ActiveVendorSource] != "replacement.go" || len(resp.OverlayReplace) != 1 { t.Fatalf("stale bridge rewrote: %#v", resp.OverlayReplace) }
}
```
