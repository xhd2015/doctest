## Expected

- Run: slow only; no skips.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requirePaths(t, resp.RunPaths, []string{"slow"}, "run")
	if len(resp.SkippedPaths) != 0 {
		t.Fatalf("skipped=%v want empty", resp.SkippedPaths)
	}
}
```
