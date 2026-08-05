## Expected

- Non-zero exit code.
- Combined output states that `--dry-run` requires `--clean` (clear message).
- Not merely `unknown command: cache`.

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
	out := combinedOut(resp)
	if strings.Contains(out, "unknown command: cache") {
		t.Fatalf("cache not implemented yet:\n%s", out)
	}
	lower := strings.ToLower(out)
	// Prefer locked wording family: "--dry-run requires --clean"
	hasRequires := strings.Contains(lower, "requires") || strings.Contains(lower, "require")
	hasDry := strings.Contains(out, "--dry-run") || strings.Contains(lower, "dry-run")
	hasClean := strings.Contains(out, "--clean") || strings.Contains(lower, "clean")
	if !(hasRequires && hasDry && hasClean) {
		t.Fatalf("expected dry-run-requires-clean message; got:\n%s", out)
	}
}
```
