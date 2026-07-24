## Expected

- Empty run; one skip.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.RunPaths) != 0 {
		t.Fatalf("run=%v want empty", resp.RunPaths)
	}
	requirePaths(t, resp.SkippedPaths, []string{"labeled_leaf"}, "skipped")
}
```
