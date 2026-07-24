---
label: heavy
---

## Expected

- Both keys are lowercase hex.
- `Key` (go1.25.0) differs from `Key2` (go1.24.0).

## Errors

- No error.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.GoVersion == req.GoVersionB {
		t.Fatalf("precondition: GoVersion strings must differ")
	}
	if !hexKey(resp.Key) || !hexKey(resp.Key2) {
		t.Fatalf("keys must be hex: %q / %q", resp.Key, resp.Key2)
	}
	if resp.Key == resp.Key2 {
		t.Fatalf("different GoVersion must yield different keys; both = %q", resp.Key)
	}
}
```
