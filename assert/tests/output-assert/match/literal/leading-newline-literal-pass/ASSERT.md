---
label: heavy
---

## Expected
- Match succeeds because the template and actual output both start with the
  same leading newline.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```
