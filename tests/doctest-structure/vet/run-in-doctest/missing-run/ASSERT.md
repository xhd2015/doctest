## Expected

- `doctest vet` fails with a non-zero exit code.
- Output indicates `func Run` is required in `DOCTEST.md`.

## Errors

- stderr or stdout references `Run`.

## Exit Code

- Non-zero.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertVetFail(t, resp, err, "Run")
}
```