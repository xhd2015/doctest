---
label: heavy
---

## Expected

- Command succeeds with exit code 0.
- stdout contains the literal version `0.0.2`.
- stdout does not contain the unresolved placeholder `__DOCTEST_VERSION__`.

## Exit Code

- Zero.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertSkillShowVersion(t, resp, err)
}
```