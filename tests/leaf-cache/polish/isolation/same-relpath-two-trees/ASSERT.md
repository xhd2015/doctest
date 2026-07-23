---
label: heavy
---

## Expected

- No harness error.
- Both keys are lowercase hex digests.
- Keys **differ** despite identical relative spine content (tree identity).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Err != "" {
		t.Fatalf("resp.Err: %s", resp.Err)
	}
	if !hexKey(resp.Key) || !hexKey(resp.Key2) {
		t.Fatalf("expected hex keys; Key=%q Key2=%q", resp.Key, resp.Key2)
	}
	if resp.Key == resp.Key2 {
		t.Fatalf("identical relpath under different TreeRoots must not share key; got %s", resp.Key)
	}
}
```
