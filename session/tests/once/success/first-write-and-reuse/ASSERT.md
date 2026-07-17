## Expected

- First Once succeeds with non-nil JSON equal to the payload.
- Second Once succeeds with **byte-identical** raw value.
- `fn` runs exactly once.
- No error on either call.

## Side Effects

- A `value` file under the once dir holds the JSON (inspected via equality of
  returned messages; layout covered in `cache-dir-layout`).

```go
import (
	"bytes"
	"encoding/json"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("first Once error: %v", resp.Err)
	}
	if resp.SecondErr != nil {
		t.Fatalf("second Once error: %v", resp.SecondErr)
	}
	if !json.Valid(resp.Value) {
		t.Fatalf("first value is not valid JSON: %s", resp.Value)
	}
	if !bytes.Equal(resp.Value, resp.SecondValue) {
		t.Fatalf("raw values differ:\n  first:  %s\n  second: %s", resp.Value, resp.SecondValue)
	}
	want := []byte(req.JSONPayload)
	if !bytes.Equal(resp.Value, want) {
		t.Fatalf("value %s want %s", resp.Value, want)
	}
	if resp.FnCalls != 1 {
		t.Fatalf("fn calls=%d want 1", resp.FnCalls)
	}
}
```
