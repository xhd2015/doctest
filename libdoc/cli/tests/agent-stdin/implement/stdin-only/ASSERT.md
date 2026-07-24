## Expected
- `err` is non-nil but does NOT contain `"requires <prompt>"`.
- The error comes from session resolution (not being inside opencode), not from empty prompt.

## Errors
- `err` must NOT be nil and must NOT contain `"requires <prompt>"`.

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
        t.Fatalf("stdin was not read as prompt, got error: %v", resp.Err)
    }
}
```
