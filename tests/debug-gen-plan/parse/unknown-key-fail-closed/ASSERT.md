## Expected

- Parse fails with an unknown-key error.
- Error mentions `not-a-key` (or "unknown key").

## Errors

- Non-empty ParseErr; fail-closed preserved.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ParseErr == "" {
		t.Fatal("expected parse error for unknown key, got nil")
	}
	low := strings.ToLower(resp.ParseErr)
	if !strings.Contains(low, "unknown") && !strings.Contains(resp.ParseErr, "not-a-key") {
		t.Fatalf("expected unknown-key error, got: %s", resp.ParseErr)
	}
	if !strings.Contains(resp.ParseErr, "not-a-key") {
		t.Fatalf("expected error to name not-a-key, got: %s", resp.ParseErr)
	}
}
```
