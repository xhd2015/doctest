## Expected
- Parse succeeds with `RegexLine` and `InlineRegex`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummaryContains(t, resp, "RegexLine", "InlineRegex")
}
```
