## Expected
- Parse succeeds.
- Summary mentions ID placeholder (regex custom subpattern).

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireParseOK(t, resp)
	requireSummaryContains(t, resp, "ID")
}
```
