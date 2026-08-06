## Expected
- Parse succeeds: literal + `BlockOptional` + literal.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummary(t, resp, "LiteralLine+BlockOptional{}+LiteralLine")
}
```
