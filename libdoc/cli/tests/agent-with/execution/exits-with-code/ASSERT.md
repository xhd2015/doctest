## Expected
- The child program exits 42; the **doctest process** reports that failure (stderr
  mentions exit status 42). Parent process exit is non-zero (CLI main exits 1 on
  any error — it does not re-emit the child's numeric code as its own exit code).

## Exit Code
- Non-zero (doctest binary).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if resp.Err == nil {
		t.Fatal("expected error (child exit 42), got nil")
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero doctest exit, got 0\nstderr:\n%s", resp.Stderr)
	}
	// Real CLI: main always os.Exit(1) on error; child code appears in stderr.
	if !strings.Contains(resp.Stderr, "42") && !strings.Contains(resp.Err.Error(), "42") {
		t.Fatalf("expected child exit 42 reflected in stderr/err, got exit=%d err=%v stderr:\n%s",
			resp.ExitCode, resp.Err, resp.Stderr)
	}
}
```
