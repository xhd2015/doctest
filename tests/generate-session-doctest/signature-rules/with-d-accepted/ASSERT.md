## Expected

- Parsing Setup, Run (DOCTEST), and Assert documents with `d *session.Doctest` succeeds.
- `rules.CheckSetupSignature` / `CheckRunSignature` / `CheckAssertSignature` accept with-d params.
- `resp.ParseErr` is empty.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.ParseErr != "" {
		t.Fatalf("expected with-d signatures accepted, got parse/rules error: %s", resp.ParseErr)
	}
}
```
