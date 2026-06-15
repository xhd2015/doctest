## Expected
- `err` is non-nil but does NOT contain `"requires <prompt>"`.
- The combined prompt (requirement + separator + stdin) was used.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error (session resolution), got nil")
    }
    if strings.Contains(resp.Err.Error(), "requires <prompt>") {
        t.Fatalf("combined requirement+stdin not used, got error: %v", resp.Err)
    }
}
```
