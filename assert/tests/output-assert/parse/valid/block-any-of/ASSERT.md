## Expected
- Parse succeeds with `AnyOfBlock` containing two branches.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "AnyOfBlock", "branches:2")
}
```
