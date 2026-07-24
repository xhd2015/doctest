---
label: heavy
---

## Expected

- `doctest vet` fails with a non-zero exit code.
- Output mentions the missing version section.

## Errors

- stderr or stdout contains a message referencing `version`.

## Exit Code

- Non-zero.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertVetFail(t, resp, err, "version")
}
```