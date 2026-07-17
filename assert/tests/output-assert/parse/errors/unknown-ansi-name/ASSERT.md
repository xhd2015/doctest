---
label: heavy
---

## Expected
- Parse fails — unknown name without `#` prefix.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireParseErrorContains(t, resp, "orange")
}
```
