---
label: heavy
---

## Expected
- Match fails because the template requires a trailing blank line that actual
  output does not contain.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchError(t, resp)
}
```
