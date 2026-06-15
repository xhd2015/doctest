## Expected
- `err` is non-nil but does NOT contain `"requires <prompt>"`.

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
        t.Fatalf("positional arg not used as prompt, got error: %v", resp.Err)
    }
}
```
