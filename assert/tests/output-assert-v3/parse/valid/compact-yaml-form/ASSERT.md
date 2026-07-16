## Expected
- Parse succeeds.
- Summary mentions PORT number placeholder (metadata may appear).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "PORT", "number")
}
```
