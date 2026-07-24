## Expected

- Run e2e/fast_child; skip e2e/labeled_child.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requirePaths(t, resp.RunPaths, []string{"e2e/fast_child"}, "run")
	requirePaths(t, resp.SkippedPaths, []string{"e2e/labeled_child"}, "skipped")
}
```
