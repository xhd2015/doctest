## Expected

- No error.
- Exactly one `-overlay=` in `GoFlags`.
- Final `OverlayReplace` equals the user seed map (same key → same value).
- `OverlayFile` non-empty when a flag is emitted.

## Errors

- Classic TDD **RED** until materialize helper writes seed and emits one flag.

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
		t.Fatalf("unexpected error exit=%d err=%s", resp.ExitCode, errText(resp))
	}
	if resp.OverlayFlagN != 1 {
		t.Fatalf("want exactly one -overlay= flag, got n=%d flags=%v", resp.OverlayFlagN, resp.GoFlags)
	}
	if len(resp.GoFlags) != 1 || !strings.HasPrefix(resp.GoFlags[0], "-overlay=") {
		t.Fatalf("GoFlags=%v want single -overlay=…", resp.GoFlags)
	}
	if resp.OverlayFile == "" || !strings.Contains(resp.GoFlags[0], resp.OverlayFile) {
		t.Fatalf("flag path must match OverlayFile: flag=%v file=%q", resp.GoFlags, resp.OverlayFile)
	}
	for k, want := range req.UserReplace {
		if got := resp.OverlayReplace[k]; got != want {
			t.Fatalf("Replace[%q]=%q want %q; full=%v", k, got, want, resp.OverlayReplace)
		}
	}
	if len(resp.OverlayReplace) != len(req.UserReplace) {
		t.Fatalf("Replace size=%d want %d map=%v", len(resp.OverlayReplace), len(req.UserReplace), resp.OverlayReplace)
	}
}
```
