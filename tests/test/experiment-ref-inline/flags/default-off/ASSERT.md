## Expected

- Parse succeeds.
- `opts.ExperimentRefInsteadOfInline == false`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ParseErr != "" {
		t.Fatalf("parse error: %s", resp.ParseErr)
	}
	if resp.Opts.ExperimentRefInsteadOfInline {
		t.Fatal("ExperimentRefInsteadOfInline=true by default; want false (flag off by default)")
	}
}
```
