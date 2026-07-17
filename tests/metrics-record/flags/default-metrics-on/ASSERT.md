## Expected

- Parse succeeds.
- `opts.NoMetrics == false`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ParseErr != "" {
		t.Fatalf("parse error: %s", resp.ParseErr)
	}
	if resp.Opts.NoMetrics {
		t.Fatal("NoMetrics=true by default; want false (metrics on)")
	}
}
```
