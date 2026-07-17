## Expected

- First Once returns an error whose message is `boom` (or contains `boom`).
- Second Once returns an error (same or equivalent message).
- `fn` runs exactly once (error persistence — no retry).
- No success value on either call.

## Errors

- Error path is intentional; harness Run does not fail.

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
		t.Fatal("expected first Once error")
	}
	if !strings.Contains(resp.Err.Error(), "boom") {
		t.Fatalf("first err=%v want boom", resp.Err)
	}
	if resp.SecondErr == nil {
		t.Fatal("expected second Once error")
	}
	if !strings.Contains(resp.SecondErr.Error(), "boom") {
		t.Fatalf("second err=%v want boom", resp.SecondErr)
	}
	if resp.FnCalls != 1 {
		t.Fatalf("fn calls=%d want 1 (error must be persisted)", resp.FnCalls)
	}
	if len(resp.Value) != 0 {
		t.Fatalf("expected no success value, got %s", resp.Value)
	}
}
```
