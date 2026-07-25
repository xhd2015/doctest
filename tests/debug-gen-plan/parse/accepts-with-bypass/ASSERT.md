## Expected

- Parse accepts both keys.
- GenPlan true and BypassGoTest true.

## Errors

- No parse error.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ParseErr != "" {
		t.Fatalf("expected combined parse to succeed, got: %s", resp.ParseErr)
	}
	if !resp.GenPlan {
		t.Fatal("expected GenPlan true")
	}
	if !resp.BypassGoTest {
		t.Fatal("expected BypassGoTest true")
	}
}
```
