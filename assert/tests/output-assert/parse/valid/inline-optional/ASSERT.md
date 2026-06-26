## Expected
- Parse succeeds with `InlineOptional` segment.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "InlineOptional")
}
```
