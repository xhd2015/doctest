---
label: heavy
---

## Expected
- Parse succeeds with `AnsiColor` spec `gray`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "AnsiColor", "gray")
}
```
