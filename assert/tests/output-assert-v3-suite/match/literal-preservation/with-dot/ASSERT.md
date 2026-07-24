---
label: heavy
---

## Expected
- Match succeeds — escaped `\.` matches a literal dot only.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```
