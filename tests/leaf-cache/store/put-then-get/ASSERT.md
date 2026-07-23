---
label: heavy
---

## Expected

- No error.
- After PutPass, GetPass returns true (`Hit == true`).
- Response Key matches the request key.

## Side Effects

- Pass marker exists only under the test StoreRoot.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Key != req.Key {
		t.Fatalf("Key = %q, want %q", resp.Key, req.Key)
	}
	if !resp.Hit {
		t.Fatal("GetPass after PutPass must be true")
	}
}
```
