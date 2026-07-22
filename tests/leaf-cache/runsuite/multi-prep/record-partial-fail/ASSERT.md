## Expected

- No harness error.
- After RecordPasses: `Hit` (treeA) is **false** — failed identity not stored.
- `HitB` (treeB) is **true** — non-failed identity was PutPass'd.

## Errors

- RecordPasses does not return an error (best-effort store I/O; test store is writable).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Hit {
		t.Fatalf("failed treeA must not PutPass; Hit=true keyA=%q", resp.Key)
	}
	if !resp.HitB {
		t.Fatalf("non-failed treeB must PutPass; HitB=false keyB=%q", resp.Key2)
	}
}
```
