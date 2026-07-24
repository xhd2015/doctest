## Expected

- Parse fails (non-zero exit mapping).
- Error reports the unknown runner option.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected nonzero exit")
	}
	msg := resp.Stderr + resp.ParseErr
	if !strings.Contains(msg, "unknown flag") &&
		!strings.Contains(msg, "unknown option") &&
		!strings.Contains(msg, "unrecognized flag") {
		t.Fatalf("stderr/parse missing unknown flag message:\n%s", msg)
	}
}
```
