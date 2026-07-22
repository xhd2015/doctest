## Expected

- Run paths: both, slow (sorted).
- Three skips.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requirePaths(t, resp.RunPaths, []string{"both", "slow"}, "run")
	if len(resp.SkippedPaths) != 3 {
		t.Fatalf("skipped=%v want 3", resp.SkippedPaths)
	}
}
```
