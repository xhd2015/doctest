## Expected
- `err` is non-nil and contains `"requires <prompt>"`.
- Empty pipe means no prompt available (stdin EOF with no data = empty prompt).

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error for empty stdin pipe, got nil")
    }
    if !strings.Contains(resp.Err.Error(), "requires <prompt>") {
        t.Fatalf("expected 'requires <prompt>' for empty pipe, got: %v", resp.Err)
    }
}
```
