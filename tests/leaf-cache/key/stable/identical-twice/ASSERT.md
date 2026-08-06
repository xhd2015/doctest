## Expected

- `ComputeLeafKey` succeeds twice with no error.
- `Key` and `Key2` are equal.
- Both are non-empty lowercase hex strings of reasonable length (≥16 chars).

## Errors

- No error from either call.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if !hexKey(resp.Key) {
		t.Fatalf("Key %q is not a lowercase hex digest", resp.Key)
	}
	if !hexKey(resp.Key2) {
		t.Fatalf("Key2 %q is not a lowercase hex digest", resp.Key2)
	}
	if resp.Key != resp.Key2 {
		t.Fatalf("keys differ on identical inputs: %q vs %q", resp.Key, resp.Key2)
	}
}
```
