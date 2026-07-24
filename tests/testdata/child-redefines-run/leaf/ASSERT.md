## Expected
- The root Run should be used, returning a Result prefixed with "root:".
- If the child Run were mistakenly used, Result would start with "child:".

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !strings.HasPrefix(resp.Result, "root:") {
        t.Fatalf("expected root Run result, got %q", resp.Result)
    }
}
```
