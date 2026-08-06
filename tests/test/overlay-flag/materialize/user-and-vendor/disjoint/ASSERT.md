## Expected

- No error; exactly one `-overlay=` flag.
- User package mapping present; vendor go.mod → placeholder present.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if resp.ExitCode != 0 || resp.ErrMsg != "" {
		t.Fatalf("unexpected error: %s", errText(resp))
	}
	if resp.OverlayFlagN != 1 || len(resp.GoFlags) != 1 || !strings.HasPrefix(resp.GoFlags[0], "-overlay=") {
		t.Fatalf("want single -overlay=, got n=%d flags=%v", resp.OverlayFlagN, resp.GoFlags)
	}
	if got := resp.OverlayReplace[resp.ProjectSource]; got != "from-user" {
		t.Fatalf("user package key: got %q map=%v", got, resp.OverlayReplace)
	}
	if got := resp.OverlayReplace[resp.ActiveVendorGoMod]; got != resp.ActiveBridgeSource {
		t.Fatalf("vendor go.mod key: got %q want %q map=%v", got, resp.ActiveBridgeSource, resp.OverlayReplace)
	}
}
```
