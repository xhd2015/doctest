---
label: heavy
---

## Expected

- Keys are lowercase hex.
- Changing the replace target's `go.mod` changes the leaf key.

## Errors

- No error.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hexKey(resp.Key) || !hexKey(resp.Key2) {
		t.Fatalf("keys must be hex: %q / %q", resp.Key, resp.Key2)
	}
	if resp.Key == resp.Key2 {
		t.Fatalf("replace lib go.mod change must alter key; both = %q", resp.Key)
	}
}
```
