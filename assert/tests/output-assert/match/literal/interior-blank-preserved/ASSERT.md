---
label: heavy
---

## Expected
- Match succeeds. Strict parsing preserves the interior empty line between
  `a` and `b`, and it matches the actual.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchOK(t, resp)
}
```
