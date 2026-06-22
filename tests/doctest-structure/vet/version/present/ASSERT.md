## Expected

- `doctest vet` succeeds with exit code 0.

## Exit Code

- Zero.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertVetPass(t, resp, err)
}
```