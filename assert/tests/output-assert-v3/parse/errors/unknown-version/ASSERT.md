## Expected
- Parse fails for unknown version (must not silently fall back to v1).

## Errors
- Error mentions version or unknown.

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
	if !strings.Contains(msg, "version") && !strings.Contains(msg, "unknown") && !strings.Contains(msg, "9") {
		t.Fatalf("parse error should mention unknown version: %q", resp.ParseError)
	}
}
```
