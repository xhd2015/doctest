## Expected

- `doctest test` succeeds with exit code 0.
- The runnable leaf passes.

## Exit Code

- Zero.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertCommandPass(t, resp, err)
}
```