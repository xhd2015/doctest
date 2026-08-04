## Expected

- The executor writes the pre-created file.
- One, and only one, standard Go overlay flag is returned.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil || resp.ErrMsg != "" { t.Fatalf("unexpected error: %v %s", err, resp.ErrMsg) }
	if !resp.OverlayDirExists || !resp.OverlayFileExists || resp.OverlayFileSize == 0 { t.Fatalf("expected populated shared overlay: %#v", resp) }
	want := "-overlay=" + resp.OverlayFile
	if len(resp.GoFlags) != 1 || resp.GoFlags[0] != want { t.Fatalf("flags=%#v want %#v", resp.GoFlags, []string{want}) }
}
```
