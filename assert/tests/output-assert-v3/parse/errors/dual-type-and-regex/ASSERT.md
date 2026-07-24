## Expected
- Parse fails because `type=` and `regex=` cannot both be set on one placeholder.

## Errors
- Error mentions type, regex, or that both cannot be combined.

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
	if !strings.Contains(msg, "type") && !strings.Contains(msg, "regex") {
		t.Fatalf("parse error should mention type/regex dual def: %q", resp.ParseError)
	}
}
```
