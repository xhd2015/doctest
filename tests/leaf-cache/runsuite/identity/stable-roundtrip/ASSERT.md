## Expected

- No harness error.
- Two formats of the same `(TreeRoot, "leaf")` are equal (`Key == Key2`).
- Different leaf rel under the same tree yields different identity
  (`Identity != Identity2`).
- All returned identity strings are non-empty.

## Errors

- No error from FormatLeafIdentity.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Key == "" || resp.Key2 == "" {
		t.Fatalf("stable formats must be non-empty: %q / %q", resp.Key, resp.Key2)
	}
	if resp.Key != resp.Key2 {
		t.Fatalf("FormatLeafIdentity must be stable: first=%q second=%q", resp.Key, resp.Key2)
	}
	if resp.Identity == "" || resp.Identity2 == "" {
		t.Fatalf("rel-distinct identities must be non-empty: %q / %q", resp.Identity, resp.Identity2)
	}
	if resp.Identity == resp.Identity2 {
		t.Fatalf("different leaf rels under same tree must differ: both=%q", resp.Identity)
	}
	// Key is Format(tree,"leaf"); Identity is same for leaf; Identity2 is "other".
	if resp.Key != resp.Identity {
		t.Fatalf("Key (stable leaf) should equal Identity (leaf): Key=%q Identity=%q", resp.Key, resp.Identity)
	}
}
```
