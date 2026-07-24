---
label: heavy
---

## Expected
- Parse succeeds.
- Summary mentions PORT placeholder with number type.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "PORT", "number")
}
```