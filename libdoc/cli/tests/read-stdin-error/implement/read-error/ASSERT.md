## Expected
- `resp.Err` is non-nil — the `io.ReadAll()` error propagates through `cli.Run()`.
- The error message indicates reading a directory is not allowed.

## Side Effects
- None (no implementation is invoked because the error occurs before `implementer.Run`).

## Errors
- Error propagated from `readStdinIfPresent()` via `runAgentImplement()`.

## Exit Code
- N/A (test calls `cli.Run()` directly)

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err == nil {
		t.Fatal("expected error from directory stdin, got nil")
	}
	errStr := resp.Err.Error()
	if !strings.Contains(errStr, "is a directory") &&
		!strings.Contains(errStr, "read") {
		t.Fatalf("expected read error (is a directory), got: %v", resp.Err)
	}
}
```
