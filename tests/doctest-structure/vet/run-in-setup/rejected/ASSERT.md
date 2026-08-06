## Expected

- `doctest vet` fails with a non-zero exit code.
- Output indicates types must be defined in `DOCTEST.md`, not root `SETUP.md`.

## Errors

- stderr or stdout references `DOCTEST.md`.

## Exit Code

- Non-zero.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertVetFail(t, resp, err, "DOCTEST.md")
}
```