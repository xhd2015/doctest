## Expected
- Match succeeds. The leading empty template line (from the raw-string `\n`)
  is stripped before matching, so the pattern is contains-only and the
  substring matches mid-line.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```
