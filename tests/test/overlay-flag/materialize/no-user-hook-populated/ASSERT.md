## Expected

- No error; exactly one `-overlay=` flag.
- Replace contains hook mapping for project-source; no requirement for user keys.

## Errors

- Classic TDD **RED** until empty-user path on `ApplyPreTestHooksWithUserOverlay`
  matches today's hook-populated behavior (may be GREEN once helper delegates to
  existing ApplyPreTestHooksWithVendorBridges).

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
	if got := resp.OverlayReplace[resp.ProjectSource]; got != "from-hook" {
		t.Fatalf("hook mapping missing: got %q map=%v", got, resp.OverlayReplace)
	}
}
```
