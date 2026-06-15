## Expected
- `err` is non-nil and contains `"read requirement file"`.
- Missing requirement file produces an error before the prompt check.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error for missing requirement file, got nil")
    }
    if !strings.Contains(resp.Err.Error(), "read requirement file") {
        t.Fatalf("expected 'read requirement file' error, got: %v", resp.Err)
    }
}
```
