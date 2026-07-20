## Expected

- No error.
- After PutPass on root A: `Hit == true` (store A).
- Store B never saw PutPass: `HitB == false`.

## Side Effects

- Only files under StoreRoot (A) are created; StoreRootB stays without that key.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Hit {
		t.Fatal("store A must GetPass true after PutPass")
	}
	if resp.HitB {
		t.Fatal("store B must not see PutPass from store A (root isolation)")
	}
}
```
