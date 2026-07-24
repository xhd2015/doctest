## Expected
- Parse succeeds on the v3 path (not legacy_v1 literal fallback).
- Summary mentions USER string placeholder.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "USER", "string")
}
```
