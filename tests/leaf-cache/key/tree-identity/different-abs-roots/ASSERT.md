## Expected

- Both keys are lowercase hex digests.
- Keys **differ** despite identical relative spine content (tree identity).

## Errors

- No error from either ComputeLeafKey call.

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
		t.Fatalf("identical content under different TreeRoot must not share a key; both=%q\ntreeA=%s\ntreeB=%s",
			resp.Key, req.TreeRoot, req.TreeRootB)
	}
}
```
