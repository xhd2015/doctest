## Expected

```go
import (
	"testing"
	"github.com/xhd2015/doctest/session"
)
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil || resp.ErrMsg != "" { t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg) }
	if len(resp.Calls) != 2 || resp.Calls[0][0] != "first" || resp.Calls[1][0] != "second" { t.Fatalf("hook order=%#v", resp.Calls) }
	if resp.OverlayReplace[resp.ProjectSource] != "project.go" || resp.OverlayReplace[resp.ActiveBridgeSource] != "active.go" { t.Fatalf("final overlay=%#v", resp.OverlayReplace) }
	if _, ok := resp.OverlayReplace[resp.ActiveVendorSource]; ok { t.Fatalf("normalization did not see final hook output: %#v", resp.OverlayReplace) }
}
```
