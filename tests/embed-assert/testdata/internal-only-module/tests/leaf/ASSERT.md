## Test

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil { t.Fatal(err) }
    if resp.Message != "hi" { t.Fatalf("expected hi, got %q", resp.Message) }
}
```