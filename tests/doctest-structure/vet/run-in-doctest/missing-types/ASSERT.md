---
label: heavy
---

## Expected

- `doctest vet` fails with a non-zero exit code.
- Output indicates `Request` and `Response` are required in `DOCTEST.md`.

## Errors

- stderr or stdout references `Request`.

## Exit Code

- Non-zero.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertVetFail(t, resp, err, "Request")
}
```