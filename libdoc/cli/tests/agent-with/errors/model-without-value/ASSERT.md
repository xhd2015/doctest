## Expected
- `resp.Err` is non-nil and contains `"--model requires a value"`.

## Errors
- The parser must reject an empty `--model` value.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error for missing --model value, got nil")
    }
    if !strings.Contains(resp.Err.Error(), "--model requires a value") {
        t.Fatalf("expected error '--model requires a value', got: %v", resp.Err)
    }
}
```
