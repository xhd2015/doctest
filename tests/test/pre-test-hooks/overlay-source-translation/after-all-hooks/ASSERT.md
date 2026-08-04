## Expected

- Hooks run in order (first, second).
- Final overlay keeps project + package keys from hooks and merges go.mod pair last.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil || resp.ErrMsg != "" {
		t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg)
	}
	if len(resp.Calls) != 2 || resp.Calls[0][0] != "first" || resp.Calls[1][0] != "second" {
		t.Fatalf("hook order=%#v", resp.Calls)
	}
	if resp.OverlayReplace[resp.ProjectSource] != "project.go" ||
		resp.OverlayReplace[resp.ActiveVendorSource] != "active.go" {
		t.Fatalf("hook mappings missing: %#v", resp.OverlayReplace)
	}
	if got := resp.OverlayReplace[resp.ActiveVendorGoMod]; got != resp.ActiveBridgeSource {
		t.Fatalf("post-hook go.mod merge: got=%q want=%q overlay=%#v", got, resp.ActiveBridgeSource, resp.OverlayReplace)
	}
}
```
