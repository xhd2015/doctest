## Expected

- Non-zero exit, parse error on stderr, no PASS/FAIL summary.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected nonzero exit\nstderr:\n%s", resp.Stderr)
	}
	msg := resp.Stderr + resp.ParseErr
	if !strings.Contains(msg, "parse") && !strings.Contains(msg, "syntax") {
		t.Fatalf("stderr must report expression parse/syntax error, not generic flag error:\n%s", msg)
	}
	// In-process parse never runs the suite — no PASS/FAIL summary.
	if strings.Contains(resp.Stdout, "PASS (") || strings.Contains(resp.Stdout, "FAIL (") {
		t.Fatalf("expected no PASS/FAIL summary\nstdout:\n%s", resp.Stdout)
	}
}
```
