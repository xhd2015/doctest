## Expected
- Parse succeeds as v2 despite leading blank lines before `---`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireSummary(t, resp, "LiteralLine")
}
```