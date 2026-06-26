## Expected
- Parse fails with position information.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireParseError(t, resp)
}
```
