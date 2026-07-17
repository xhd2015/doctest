---
label: heavy
---

## Expected
- Parse fails due to label mismatch.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireParseErrorContains(t, resp, "label", "id", "wrong")
}
```
