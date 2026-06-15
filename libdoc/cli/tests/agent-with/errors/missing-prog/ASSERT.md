## Expected
- `resp.Err` is non-nil and contains `"agent with requires <prog>"`.

## Errors
- The subcommand requires at least one positional argument (the program to execute).

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.Err == nil {
        t.Fatal("expected error for missing <prog>, got nil")
    }
    if !strings.Contains(resp.Err.Error(), "agent with requires <prog>") {
        t.Fatalf("expected error 'agent with requires <prog>', got: %v", resp.Err)
    }
}
```
