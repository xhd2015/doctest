## Expected
- `resp.Err` is non-nil — the `os.Stdin.Stat()` error propagates through `cli.Run()`.
- The error message indicates a closed/bad file descriptor.

## Side Effects
- None (no implementation is invoked because the error occurs before `designer.Run`).

## Errors
- Error propagated from `readStdinIfPresent()` via `runAgentDesign()`.

## Exit Code
- N/A (test calls `cli.Run()` directly)

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.Err == nil {
		t.Fatal("expected error from closed stdin, got nil")
	}
	errStr := resp.Err.Error()
	if !strings.Contains(errStr, "file already closed") &&
		!strings.Contains(errStr, "bad file descriptor") {
		t.Fatalf("expected stat error (closed/bad fd), got: %v", resp.Err)
	}
}
```
