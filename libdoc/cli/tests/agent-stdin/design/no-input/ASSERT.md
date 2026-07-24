## Expected
- `err` is non-nil and contains `"requires <prompt>"`.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error for no-input case, got nil")
    }
    if !strings.Contains(resp.Err.Error(), "requires <prompt>") {
        t.Fatalf("expected error 'requires <prompt>', got: %v", resp.Err)
    }
}
```
