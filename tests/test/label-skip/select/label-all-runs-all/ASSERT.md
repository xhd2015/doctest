## Expected

- Run both leaves; no skips.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requirePaths(t, resp.RunPaths, []string{"fast_leaf", "labeled_leaf"}, "run")
	if len(resp.SkippedPaths) != 0 {
		t.Fatalf("skipped=%v", resp.SkippedPaths)
	}
}
```
