---
label: heavy
---

## Expected
- Parse succeeds with `PatternLine[Hint:id]`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "PatternLine", "Hint:id")
}
```
