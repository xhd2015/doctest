## Expected

- Non-zero exit.
- stderr indicates `--color` and `--no-color` cannot be specified together.

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
	se := strings.ToLower(resp.Stderr)
	ok := (strings.Contains(se, "color") && strings.Contains(se, "no-color")) ||
		strings.Contains(se, "together") ||
		strings.Contains(se, "conflict") ||
		strings.Contains(se, "mutually")
	if !ok {
		t.Fatalf("expected color flag conflict message:\n%s", resp.Stderr)
	}
}
```
