## Expected

- No harness error.
- `SkipPaths` is empty (or nil).
- Identities for both trees are still non-empty (plan still maps keys).
- Optional: `Key` and `Key2` store hex keys non-empty when returned.

## Errors

- PreparePassPlan returns nil error.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.SkipPaths) != 0 {
		t.Fatalf("skip disabled → empty Skip; got %v", resp.SkipPaths)
	}
	if resp.Identity == "" || resp.Identity2 == "" {
		t.Fatalf("keys/identities still required when skip disabled: %q / %q", resp.Identity, resp.Identity2)
	}
	if resp.Key == "" || resp.Key2 == "" {
		t.Fatalf("store keys should still be computed: %q / %q", resp.Key, resp.Key2)
	}
}
```
