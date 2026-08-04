## Expected

- The response reports the hook failure with the hook command context.
- The second hook does not run.
- No Go overlay argument is returned.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil { t.Fatalf("Run infrastructure error: %v", err) }
	if len(resp.Calls) != 1 || resp.Calls[0][0] != "first" { t.Fatalf("hooks after failure=%#v", resp.Calls) }
	if !strings.Contains(resp.ErrMsg, "hook failed") || !strings.Contains(resp.ErrMsg, "first") { t.Fatalf("error should include hook and stderr context: %q", resp.ErrMsg) }
	if len(resp.GoFlags) != 0 { t.Fatalf("failed hook must not activate overlay: %#v", resp.GoFlags) }
}
```
