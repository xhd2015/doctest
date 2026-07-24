---
label: heavy
---

## Expected
- Parse fails for unknown placeholder type.

## Errors
- Error mentions unknown or unsupported type.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireParseErrorContains(t, resp, "boolean")
}
```