## Expected

- Exit code is 0: doctest generates valid Go from the shared-type named-result
  helper and the inner leaf test compiles and passes.
- Stdout/stderr must not contain `format imports failed` (the bug symptom when
  `writeFuncClosure` omits parentheses around `port int, alt int`).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run error: %v", err)
	}

	combined := resp.Stdout + "\n" + resp.Stderr
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0 (shared-type named results should assemble cleanly), got %d\nstdout: %s\nstderr: %s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if strings.Contains(combined, "format imports failed") {
		t.Fatalf("expected no format imports failure, got:\n%s", combined)
	}
}
```