---
label: heavy
---

## Expected

- No harness error; keys are lowercase hex.
- After editing leaf_a: **Hit == true** (sibling leaf_b key stable).
- **HitB == false** (leaf_a key changed).

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Err != "" {
		t.Fatalf("resp.Err: %s", resp.Err)
	}
	if !hexKey(resp.Key) || !hexKey(resp.Key2) {
		t.Fatalf("expected hex keys; Key=%q Key2=%q", resp.Key, resp.Key2)
	}
	if !resp.Hit {
		t.Fatalf("expected sibling leaf_b key stable (Hit=true)")
	}
	if resp.HitB {
		t.Fatalf("expected leaf_a key to change after ASSERT edit (HitB=false)")
	}
}
```
