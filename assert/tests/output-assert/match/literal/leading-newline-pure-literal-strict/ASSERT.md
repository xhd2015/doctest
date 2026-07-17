---
label: heavy
---

## Expected
- Match fails because the template's leading newline is part of the literal
  contract and the actual output starts directly with `foo`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchError(t, resp)
}
```
