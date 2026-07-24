## Expected
- The root Run should multiply the Value by 10.
- Value of 7 should produce Result of 70.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if resp.Result != 70 {
        t.Fatalf("expected Result=70, got %d", resp.Result)
    }
}
```
