## Expected
- Parse fails.

## Errors
- Bare tag without `hint:` prefix is rejected.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireParseErrorContains(t, resp, "id")
}
```
