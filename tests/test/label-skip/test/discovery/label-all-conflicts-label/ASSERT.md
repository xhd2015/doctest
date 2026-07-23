## Expected

- Non-zero exit; error mentions mutually exclusive flags.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + "\n" + resp.Stderr + "\n" + resp.ParseErr
	if !strings.Contains(combined, "mutually exclusive") {
		t.Fatalf("expected mutual exclusion message\nstdout:\n%s\nstderr:\n%s\nparse:\n%s",
			resp.Stdout, resp.Stderr, resp.ParseErr)
	}
}
```
