## Expected

- `session.Once` returns a non-nil error.
- Error mentions session id / `DOCTEST_SESSION_ID` (implementation wording may vary).
- `fn` is not called (`FnCalls == 0`).
- Returned value is empty.

## Errors

- Missing session id is a hard error, not a panic.

## Exit Code

- Test process exits 0 only when Assert passes (error path is the success for this leaf).

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp.Err == nil {
		t.Fatal("expected error when DOCTEST_SESSION_ID is missing")
	}
	msg := resp.Err.Error()
	if !strings.Contains(msg, "DOCTEST_SESSION_ID") && !strings.Contains(strings.ToLower(msg), "session") {
		t.Fatalf("error should mention session id, got: %v", resp.Err)
	}
	if resp.FnCalls != 0 {
		t.Fatalf("fn must not run when session id missing, calls=%d", resp.FnCalls)
	}
	if len(resp.Value) != 0 {
		t.Fatalf("expected empty value on error, got %s", resp.Value)
	}
}
```
