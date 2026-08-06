## Expected
- Parse succeeds with raw SGR params `38;5;208`.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "AnsiColor", "38;5;208")
}
```
