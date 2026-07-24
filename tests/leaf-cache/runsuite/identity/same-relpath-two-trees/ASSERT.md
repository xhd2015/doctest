## Expected

- No harness error.
- `Identity` and `Identity2` are both non-empty.
- `Identity != Identity2` despite identical relative path `leaf`.

## Errors

- No error from FormatLeafIdentity.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Identity == "" || resp.Identity2 == "" {
		t.Fatalf("identities must be non-empty: %q / %q", resp.Identity, resp.Identity2)
	}
	if resp.Identity == resp.Identity2 {
		t.Fatalf("same relpath under different trees must not share identity; both=%q\ntreeA=%s\ntreeB=%s",
			resp.Identity, req.TreeRoot, req.TreeRootB)
	}
}
```
