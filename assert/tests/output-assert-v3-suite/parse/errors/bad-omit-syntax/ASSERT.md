---
label: heavy
---

## Expected
- Parse fails for invalid omit syntax.

## Errors
- Error mentions omit marker or invalid count.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireParseErrorContains(t, resp, "omit")
}
```