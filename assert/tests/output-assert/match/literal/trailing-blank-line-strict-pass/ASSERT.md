---
label: heavy
---

## Expected
- Match succeeds because the template and actual output both include the
  trailing blank line.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```
