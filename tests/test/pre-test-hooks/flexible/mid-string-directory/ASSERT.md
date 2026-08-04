## Expected

- Containing `$GO_INSTRUMENT_OVERLAY_DIR` allocates the shared generated directory.
- Mid-string substitution preserves the `/extra` suffix after the expanded abs path.
- No overlay file or Go overlay argument is produced.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil || resp.ErrMsg != "" {
		t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg)
	}
	if !resp.OverlayDirExists || resp.OverlayDir == "" || !filepath.IsAbs(resp.OverlayDir) {
		t.Fatalf("directory not allocated: %#v", resp)
	}
	if filepath.Base(resp.OverlayDir) != "__overlay" {
		t.Fatalf("overlay dir=%q, want project-level __overlay", resp.OverlayDir)
	}
	if resp.OverlayFile != "" || len(resp.GoFlags) != 0 {
		t.Fatalf("unexpected file/flags: %#v", resp)
	}
	if len(resp.Calls) != 1 {
		t.Fatalf("calls=%#v", resp.Calls)
	}
	want := "--dir=" + resp.OverlayDir + "/extra"
	if got := resp.Calls[0][1]; got != want {
		t.Fatalf("mid-string dir arg=%q want %q", got, want)
	}
}
```
