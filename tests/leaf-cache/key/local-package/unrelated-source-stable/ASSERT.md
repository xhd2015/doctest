---
label: heavy
---

## Expected

- Keys are lowercase hex.
- Changing `unrelated/noise.go` does **not** change the leaf key (`Key == Key2`).

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
	if resp.Key != resp.Key2 {
		t.Fatalf("unrelated local package must not alter key: %q vs %q", resp.Key, resp.Key2)
	}
}
```
