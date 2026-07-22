## Expected

- No harness error; keys are lowercase hex samples.
- After editing `shared/a`: **Hit == false** (leaf-ab-1 / leaf-ab-2 keys changed).
- **HitB == true** (leaf-d key stable).

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
	if resp.Hit {
		t.Fatalf("after shared/a edit expected ab leaves keys to change (Hit=false)")
	}
	if !resp.HitB {
		t.Fatalf("after shared/a edit expected leaf-d key stable (HitB=true)")
	}
}
```
