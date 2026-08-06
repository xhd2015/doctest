## Expected
- Parse succeeds with tokens `[bold, gray]`.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "AnsiColor", "bold", "gray")
}
```
