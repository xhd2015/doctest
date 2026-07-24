## Expected
- Parse fails for invalid omit marker syntax.

## Errors
- Non-empty parse error (omit / lines / syntax related).

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireParseError(t, resp)
}
```
