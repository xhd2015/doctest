## Expected

- Run: fast_leaf; skip: labeled_leaf.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requirePaths(t, resp.RunPaths, []string{"fast_leaf"}, "run")
	requirePaths(t, resp.SkippedPaths, []string{"labeled_leaf"}, "skipped")
}
```
