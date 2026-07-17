---
label: heavy
---

## Expected
- Match succeeds for first alternation branch.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```