## Expected
- Match fails when repeated placeholders capture different strings.

## Errors
- Error mentions the placeholder name `ID` (binding mismatch).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireMatchError(t, resp, "ID")
}
```
