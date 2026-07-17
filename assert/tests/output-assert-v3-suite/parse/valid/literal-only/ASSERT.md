---
label: heavy
---

## Expected
- Parse succeeds with two content lines summarized as `RegexLine+RegexLine`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummary(t, resp, "RegexLine+RegexLine")
}
```
