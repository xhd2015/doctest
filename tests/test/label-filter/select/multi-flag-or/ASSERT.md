## Expected

- Same run set as OR expression: both, flaky, slow.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requirePaths(t, resp.RunPaths, []string{"both", "flaky", "slow"}, "run")
	requirePaths(t, resp.SkippedPaths, []string{"fast", "ui"}, "skipped")
}
```
