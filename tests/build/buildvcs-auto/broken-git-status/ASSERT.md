## Expected

- `go build -buildvcs=auto` fails (non-zero exit).
- Combined output contains `error obtaining VCS status`.
- Combined output hints at `-buildvcs=false`.

## Exit Code

- non-zero

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ExitCode == 0 {
		t.Fatalf("broken git: expected non-zero exit under -buildvcs=auto, got 0\n%s", resp.Output)
	}
	if !strings.Contains(resp.Output, "error obtaining VCS status") {
		t.Fatalf("broken git: missing 'error obtaining VCS status' in output:\n%s", resp.Output)
	}
	if !strings.Contains(resp.Output, "-buildvcs=false") {
		t.Fatalf("broken git: missing -buildvcs=false hint:\n%s", resp.Output)
	}
}
```
