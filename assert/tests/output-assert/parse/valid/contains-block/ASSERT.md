## Expected
- Parse succeeds with `ContainsBlock` fragments:3.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "ContainsBlock", "fragments:3")
}
```
