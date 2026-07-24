## Expected
- The root Run returns a Message containing "hello doctest".
- This confirms the root Run is the one executed.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.Contains(resp.Message, "hello doctest") {
        t.Fatalf("expected greeting, got %q", resp.Message)
    }
}
```
