## Expected
- Template parses as v3 (omit marker recognized).
- Match fails due to omit count mismatch.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	// Must be v3 path (not v1 literal fallback of the YAML header).
	requireSummaryContains(t, resp, "OmitLine")
	requireMatchError(t, resp)
}
```
