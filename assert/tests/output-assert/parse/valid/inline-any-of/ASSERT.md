---
label: heavy
---

## Expected
- Parse succeeds with inline any-of segments.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "InlineAnyOf", "branches:2")
}
```
