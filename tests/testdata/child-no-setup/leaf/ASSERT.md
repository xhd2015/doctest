## Expected
- This leaf would assert that Run doubles the Value if it executed.
  It exists to make the fixture a runnable tree for discovery.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    _ = req
    _ = resp
    _ = err
}
```
