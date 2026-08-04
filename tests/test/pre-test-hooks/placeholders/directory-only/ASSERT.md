## Expected

- One generated directory exists beneath the supplied generated root.
- The exact placeholder is replaced in the hook command.
- No overlay file or Go overlay argument exists.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil || resp.ErrMsg != "" { t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg) }
	if !resp.OverlayDirExists || resp.OverlayDir == "" || !filepath.IsAbs(resp.OverlayDir) { t.Fatalf("directory not allocated: %#v", resp) }
	if filepath.Base(resp.OverlayDir) != "__overlay" { t.Fatalf("overlay dir=%q, want project-level __overlay", resp.OverlayDir) }
	if resp.OverlayFile != "" || len(resp.GoFlags) != 0 { t.Fatalf("unexpected file/flags: %#v", resp) }
	if got := resp.Calls[0][2]; got != resp.OverlayDir { t.Fatalf("dir placeholder=%q want %q", got, resp.OverlayDir) }
}
```
