## Expected

- `session.Once` returns a non-nil error about empty key.
- `fn` is not called.
- Returned value is empty.

## Errors

- Empty key is rejected before disk I/O / fn.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp.Err == nil {
		t.Fatal("expected error for empty key")
	}
	msg := strings.ToLower(resp.Err.Error())
	if !strings.Contains(msg, "empty") && !strings.Contains(msg, "key") {
		t.Fatalf("error should mention empty key, got: %v", resp.Err)
	}
	if resp.FnCalls != 0 {
		t.Fatalf("fn must not run for empty key, calls=%d", resp.FnCalls)
	}
	if len(resp.Value) != 0 {
		t.Fatalf("expected empty value, got %s", resp.Value)
	}
}
```
