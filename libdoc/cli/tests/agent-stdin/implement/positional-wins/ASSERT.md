## Expected
- `err` is non-nil but does NOT contain `"requires <prompt>"`.
- Positional arg takes priority over stdin (stdin is ignored).

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error (session resolution), got nil")
    }
    if strings.Contains(resp.Err.Error(), "requires <prompt>") {
        t.Fatalf("prompt was empty despite positional args, got error: %v", resp.Err)
    }
}
```
