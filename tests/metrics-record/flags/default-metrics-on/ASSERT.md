## Expected

- Parse succeeds.
- `opts.MetricsOn == false`.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ParseErr != "" {
		t.Fatalf("parse error: %s", resp.ParseErr)
	}
	if resp.Opts.MetricsOn {
		t.Fatal("MetricsOn=true by default; want false (metrics off by default)")
	}
}
```
