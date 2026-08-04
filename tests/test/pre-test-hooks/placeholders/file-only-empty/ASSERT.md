## Expected

- One absolute `.json` overlay file exists and has size zero before/after the hook.
- The parent directory is usable even without a directory placeholder.
- No `-overlay` argument is produced.

```go
import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil || resp.ErrMsg != "" { t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg) }
	if !resp.OverlayFileExists || resp.OverlayFileSize != 0 || !filepath.IsAbs(resp.OverlayFile) || !strings.HasSuffix(resp.OverlayFile, ".json") { t.Fatalf("overlay file=%q exists=%v size=%d", resp.OverlayFile, resp.OverlayFileExists, resp.OverlayFileSize) }
	if filepath.Base(filepath.Dir(resp.OverlayFile)) != "__overlay" { t.Fatalf("overlay file=%q, want project-level __overlay/overlay.json", resp.OverlayFile) }
	if len(resp.GoFlags) != 0 { t.Fatalf("empty overlay must not add flags: %#v", resp.GoFlags) }
	if got := resp.Calls[0][2]; got != resp.OverlayFile { t.Fatalf("file placeholder=%q want %q", got, resp.OverlayFile) }
}
```
