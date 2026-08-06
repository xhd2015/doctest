## Expected
- Match succeeds. The template is written without a trailing empty line, so
  strict parsing produces a contains-only pattern and the fragment is found.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```
