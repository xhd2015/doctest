## Expected
- Parse succeeds via legacy_v1 (literal lines, not v3 AST).
- Summary shows v1 `LiteralLine` shape, not v3 `RegexLine` / placeholder nodes.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummary(t, resp, "LiteralLine×3")
	requireSummaryNotContains(t, resp, "Placeholder", "OmitLine", "RegexLine")
}
```
