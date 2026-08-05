## Expected

- Exit code 0.
- Stdout is usage for `cache` (mentions the command name).
- Usage documents `--clean` and `--dry-run`.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	requireOK(t, resp, err)
	out := combinedOut(resp)
	// Must be real cache help — not merely falling through as unknown command
	// with a generic usage dump that omits cache-specific flags.
	if strings.Contains(out, "unknown command: cache") {
		t.Fatalf("cache not implemented yet (need cache --help):\n%s", out)
	}
	mustContain(t, out, "cache", "usage")
	mustContain(t, out, "--clean", "usage")
	mustContain(t, out, "--dry-run", "usage")
}
```
