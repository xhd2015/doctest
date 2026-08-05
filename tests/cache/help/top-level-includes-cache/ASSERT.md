## Expected

- Exit code 0.
- Stdout contains top-level usage.
- Command list includes a `cache` entry as its own command line.

## Exit Code

- 0

```go
import (
	"regexp"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	requireOK(t, resp, err)
	mustContain(t, resp.Stdout, "Usage: doctest", "top-level help")
	// cache must appear as its own command entry (not only as a substring of
	// DOCTEST_CACHE_HOME or similar). RED until top-level usage lists cache.
	if !regexp.MustCompile(`(?m)^[ \t]+cache\b`).MatchString(resp.Stdout) {
		t.Fatalf("stdout missing cache command entry:\n%s", resp.Stdout)
	}
}
```
