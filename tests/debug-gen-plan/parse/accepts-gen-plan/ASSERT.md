## Expected

- `debug.Parse("gen-plan=1")` returns no error.
- `Settings.GenPlan` is true (Response.GenPlan).
- BypassGoTest remains false.

## Errors

- No parse error.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ParseErr != "" {
		t.Fatalf("expected Parse to accept gen-plan=1, got error: %s", resp.ParseErr)
	}
	if !resp.GenPlan {
		t.Fatal("expected Settings.GenPlan true after gen-plan=1")
	}
	if resp.BypassGoTest {
		t.Fatal("expected BypassGoTest false when only gen-plan=1")
	}
}
```
