## Expected
- Match succeeds. The template is written without a leading empty line, so
  strict parsing produces a contains-only pattern and the substring matches
  mid-line.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```
