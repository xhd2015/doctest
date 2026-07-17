---
label: heavy
---

## Expected
- Match fails due to omit count mismatch.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchError(t, resp)
}
```