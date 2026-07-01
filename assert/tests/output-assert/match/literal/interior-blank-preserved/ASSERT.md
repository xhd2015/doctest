## Expected
- Match succeeds. P4 strips only leading/trailing empty lines; the interior
  empty line between `a` and `b` remains meaningful and matches the actual.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```
