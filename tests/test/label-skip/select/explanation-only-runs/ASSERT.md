## Expected

- Run explained_leaf; no skips.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requirePaths(t, resp.RunPaths, []string{"explained_leaf"}, "run")
	if len(resp.SkippedPaths) != 0 {
		t.Fatalf("skipped=%v", resp.SkippedPaths)
	}
}
```
