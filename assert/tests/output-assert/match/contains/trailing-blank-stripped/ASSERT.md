## Expected
- Match succeeds. The trailing empty template line is stripped before
  matching, so the pattern is contains-only and the fragment is found.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```
