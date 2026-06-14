## Expected
- Run returns a non-nil Response with `Status == "ok"` and `Message == "processed: test"`.
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
    if resp.Status != "ok" {
        t.Fatalf("expected Status 'ok', got %q", resp.Status)
    }
    if resp.Message != "processed: test" {
        t.Fatalf("expected Message 'processed: test', got %q", resp.Message)
    }
}
```
