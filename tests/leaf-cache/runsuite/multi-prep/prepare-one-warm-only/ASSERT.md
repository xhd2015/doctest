## Expected

- No harness error.
- `SkipPaths` has exactly one entry: treeA's `Identity`.
- treeB's `Identity2` is not in `SkipPaths`.

## Errors

- PreparePassPlan returns nil error.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Identity == "" || resp.Identity2 == "" {
		t.Fatalf("expected identities: %q / %q", resp.Identity, resp.Identity2)
	}
	if len(resp.SkipPaths) != 1 {
		t.Fatalf("only treeA warm → 1 skip; got %d: %v", len(resp.SkipPaths), resp.SkipPaths)
	}
	if resp.SkipPaths[0] != resp.Identity {
		t.Fatalf("skip[0]=%q want treeA identity %q", resp.SkipPaths[0], resp.Identity)
	}
	for _, s := range resp.SkipPaths {
		if s == resp.Identity2 {
			t.Fatalf("treeB identity must not be skipped when cold: %q in %v", resp.Identity2, resp.SkipPaths)
		}
	}
}
```
