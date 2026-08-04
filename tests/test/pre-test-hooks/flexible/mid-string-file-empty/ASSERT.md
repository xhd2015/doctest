## Expected

- Containing `$GO_INSTRUMENT_OVERLAY_FILE` in any arg allocates the shared empty overlay JSON.
- The placeholder is replaced mid-string; the `--overlay=` prefix remains.
- Empty file contributes no Go overlay flag.

```go
import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil || resp.ErrMsg != "" {
		t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg)
	}
	if !resp.OverlayFileExists || resp.OverlayFileSize != 0 || !filepath.IsAbs(resp.OverlayFile) || !strings.HasSuffix(resp.OverlayFile, ".json") {
		t.Fatalf("overlay file=%q exists=%v size=%d", resp.OverlayFile, resp.OverlayFileExists, resp.OverlayFileSize)
	}
	if filepath.Base(filepath.Dir(resp.OverlayFile)) != "__overlay" {
		t.Fatalf("overlay file=%q, want project-level __overlay/overlay.json", resp.OverlayFile)
	}
	if len(resp.GoFlags) != 0 {
		t.Fatalf("empty overlay must not add flags: %#v", resp.GoFlags)
	}
	if len(resp.Calls) != 1 {
		t.Fatalf("calls=%#v", resp.Calls)
	}
	want := "--overlay=" + resp.OverlayFile
	if got := resp.Calls[0][1]; got != want {
		t.Fatalf("mid-string file arg=%q want %q", got, want)
	}
}
```
