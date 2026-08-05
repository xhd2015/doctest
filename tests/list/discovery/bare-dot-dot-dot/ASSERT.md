## Expected

- Non-zero exit.
- stderr contains `bare '...' pattern is not supported` (same family as test/vet).

## Exit Code

- non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	requireFail(t, resp, err)
	if !strings.Contains(resp.Stderr, "bare '...' pattern is not supported") {
		t.Fatalf("stderr missing bare ... message:\n%s", resp.Stderr)
	}
}
```
