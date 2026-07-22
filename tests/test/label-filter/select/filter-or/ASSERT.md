## Expected

- Run paths: both, heavy, slow.
- Two skips (fast, ui).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requirePaths(t, resp.RunPaths, []string{"both", "heavy", "slow"}, "run")
	requirePaths(t, resp.SkippedPaths, []string{"fast", "ui"}, "skipped")
}
```
