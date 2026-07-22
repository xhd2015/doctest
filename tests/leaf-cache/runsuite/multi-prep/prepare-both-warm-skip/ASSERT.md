## Expected

- No harness error.
- `SkipPaths` has length 2.
- Both `Identity` (treeA) and `Identity2` (treeB) appear in `SkipPaths`.
- The two skip entries differ (tree-qualified).

## Errors

- PreparePassPlan returns nil error.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Identity == "" || resp.Identity2 == "" {
		t.Fatalf("expected identities for both trees: %q / %q", resp.Identity, resp.Identity2)
	}
	if resp.Identity == resp.Identity2 {
		t.Fatalf("identities must be tree-qualified and distinct; both=%q", resp.Identity)
	}
	if len(resp.SkipPaths) != 2 {
		t.Fatalf("both warm → 2 skip identities; got %d: %v", len(resp.SkipPaths), resp.SkipPaths)
	}
	if !containsStr(resp.SkipPaths, resp.Identity) {
		t.Fatalf("SkipPaths missing treeA identity %q; skip=%v", resp.Identity, resp.SkipPaths)
	}
	if !containsStr(resp.SkipPaths, resp.Identity2) {
		t.Fatalf("SkipPaths missing treeB identity %q; skip=%v", resp.Identity2, resp.SkipPaths)
	}
}

func containsStr(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}
```
