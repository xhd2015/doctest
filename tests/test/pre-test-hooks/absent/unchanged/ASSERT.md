## Expected

- No overlay directory or file is allocated.
- No hook command runs.
- No Go overlay argument is contributed.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil || resp.ErrMsg != "" {
		t.Fatalf("unexpected error: run=%v hook=%s", err, resp.ErrMsg)
	}
	if len(resp.Calls) != 0 || resp.OverlayDir != "" || resp.OverlayFile != "" || len(resp.GoFlags) != 0 {
		t.Fatalf("want unchanged config application, got calls=%#v dir=%q file=%q flags=%#v", resp.Calls, resp.OverlayDir, resp.OverlayFile, resp.GoFlags)
	}
}
```
