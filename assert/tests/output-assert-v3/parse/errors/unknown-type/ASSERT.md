## Expected
- Parse fails for unknown placeholder type.

## Errors
- Error mentions type, boolean, or unknown.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireParseError(t, resp)
	msg := strings.ToLower(resp.ParseError)
	if !strings.Contains(msg, "type") && !strings.Contains(msg, "boolean") && !strings.Contains(msg, "unknown") {
		t.Fatalf("parse error should mention unknown type: %q", resp.ParseError)
	}
}
```
