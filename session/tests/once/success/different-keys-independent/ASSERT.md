## Expected

- Once for key A returns `{"key":"A"}` (or payload set in Setup).
- Once for key B returns `{"key":"B"}`.
- Values are not equal.
- `fn` runs twice (once per key).

```go
import (
	"bytes"
	"encoding/json"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("Once key A error: %v", resp.Err)
	}
	if resp.SecondErr != nil {
		t.Fatalf("Once key B error: %v", resp.SecondErr)
	}
	if !json.Valid(resp.Value) || !json.Valid(resp.SecondValue) {
		t.Fatalf("invalid JSON: %s / %s", resp.Value, resp.SecondValue)
	}
	if bytes.Equal(resp.Value, resp.SecondValue) {
		t.Fatalf("keys must be independent, both returned %s", resp.Value)
	}
	if !bytes.Equal(resp.Value, []byte(`{"key":"A"}`)) {
		t.Fatalf("key A value=%s", resp.Value)
	}
	if !bytes.Equal(resp.SecondValue, []byte(`{"key":"B"}`)) {
		t.Fatalf("key B value=%s", resp.SecondValue)
	}
	if resp.FnCalls != 2 {
		t.Fatalf("fn calls=%d want 2", resp.FnCalls)
	}
}
```
