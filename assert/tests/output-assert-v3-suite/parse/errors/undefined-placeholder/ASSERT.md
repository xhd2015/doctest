---
label: heavy
---

## Expected
- Parse fails for undefined placeholder.

## Errors
- Error mentions `__MISSING__` or undefined placeholder.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireParseErrorContains(t, resp, "MISSING")
}
```