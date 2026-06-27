## Test

```go
import (
    "testing"
    "github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil { t.Fatal(err) }
    if resp.Message != "hi" { t.Fatalf("expected hi, got %q", resp.Message) }
    assert.Output(t, "hello\n", "hello\n")
}
```