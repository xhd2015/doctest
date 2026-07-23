---
label: heavy
---

## Expected

- No harness error; keys are lowercase hex samples.
- After editing `alone/d`: **Hit == true** (leaf-ab-1 and leaf-ab-2 keys stable).
- **HitB == false** (leaf-d key changed).

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
	if !resp.Hit {
		t.Fatalf("after alone/d edit expected shared leaves (ab-1, ab-2) keys stable (Hit=true)")
	}
	if resp.HitB {
		t.Fatalf("after alone/d edit expected leaf-d key to change (HitB=false)")
	}
}
```
