## Expected
- Parse fails because the custom regex fragment does not compile.

## Errors
- Error mentions regex, compile, or invalid pattern.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	requireParseError(t, resp)
	msg := strings.ToLower(resp.ParseError)
	if !strings.Contains(msg, "regex") && !strings.Contains(msg, "compile") && !strings.Contains(msg, "invalid") && !strings.Contains(msg, "error") {
		t.Fatalf("parse error should mention invalid regex: %q", resp.ParseError)
	}
}
```
