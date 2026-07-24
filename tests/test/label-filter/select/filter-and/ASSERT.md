## Expected

- Run path: both only.
- Four skips.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requirePaths(t, resp.RunPaths, []string{"both"}, "run")
	if len(resp.SkippedPaths) != 4 {
		t.Fatalf("skipped=%v want 4", resp.SkippedPaths)
	}
}
```
