## Expected
- Match succeeds: color inner text matched literally (including `.`), PORT as number.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```
