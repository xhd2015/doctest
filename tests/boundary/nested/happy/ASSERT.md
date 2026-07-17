---
label: heavy
---

## Expected
- Run returns a non-nil Response with `Greeting == "Hello, World!"`.
- `err` is nil.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if resp == nil {
        t.Fatal("expected non-nil response")
    }
    if resp.Greeting != "Hello, World!" {
        t.Fatalf("expected Greeting 'Hello, World!', got %q", resp.Greeting)
    }
}
```
