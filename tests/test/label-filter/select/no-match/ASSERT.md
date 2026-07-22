## Expected

- No run paths; five skips.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.RunPaths) != 0 {
		t.Fatalf("run=%v want empty", resp.RunPaths)
	}
	if len(resp.SkippedPaths) != 5 {
		t.Fatalf("skipped=%v want 5", resp.SkippedPaths)
	}
}
```
