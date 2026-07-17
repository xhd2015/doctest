---
label: heavy
---

## Expected
- Parse succeeds.
- Summary captures PORT metadata including example.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "PORT", "number", "8901")
}
```