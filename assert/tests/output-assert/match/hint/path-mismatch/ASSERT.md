---
label: heavy
---

## Expected
- Match fails with `hint:path` in error.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchError(t, resp, "hint:path")
}
```
