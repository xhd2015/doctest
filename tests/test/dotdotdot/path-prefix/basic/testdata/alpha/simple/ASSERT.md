## Expected

- Pass.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatal(err)
    }
    if resp == nil {
        t.Fatal("resp is nil")
    }
}
```
