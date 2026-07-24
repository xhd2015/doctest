## Expected

- No run; skip slow with reason label filter.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.RunPaths) != 0 {
		t.Fatalf("run=%v want empty", resp.RunPaths)
	}
	requirePaths(t, resp.SkippedPaths, []string{"slow"}, "skipped")
	if len(resp.SkippedReason) != 1 || resp.SkippedReason[0] != "label filter" {
		t.Fatalf("reasons=%v", resp.SkippedReason)
	}
}
```
