## Expected

- Both keys are non-empty lowercase hex.
- After mutating parent `group/SETUP.md` Go, `Key2` differs from `Key`.

## Errors

- No error.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hexKey(resp.Key) || !hexKey(resp.Key2) {
		t.Fatalf("keys must be hex: %q / %q", resp.Key, resp.Key2)
	}
	if resp.Key == resp.Key2 {
		t.Fatalf("ancestor setup mutation must change key; both = %q", resp.Key)
	}
}
```
