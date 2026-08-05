## Expected

- Exit code 0.
- Stdout mentions Cache home and Doctest root with the injectable absolute paths.
- Indicates empty cache: 0 buckets and/or 0B / no cache wording.
- Does not claim non-zero Total for real data.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	requireOK(t, resp, err)
	out := resp.Stdout
	if strings.Contains(combinedOut(resp), "unknown command: cache") {
		t.Fatalf("cache not implemented yet:\n%s", combinedOut(resp))
	}
	// Paths from injectable Request — prove product used our CacheHome, not user cache.
	mustContain(t, out, req.CacheHome, "Cache home path")
	mustContain(t, out, req.DoctestRoot, "Doctest root path")
	// Labels (flexible casing/spacing)
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "cache home") && !strings.Contains(out, "Cache home") {
		t.Fatalf("stdout missing Cache home label:\n%s", out)
	}
	if !strings.Contains(lower, "doctest root") && !strings.Contains(out, "Doctest root") {
		t.Fatalf("stdout missing Doctest root label:\n%s", out)
	}
	// Empty signal: 0 buckets, 0B, Total: 0, or "no cache" / "empty"
	emptyOK := strings.Contains(out, "0 bucket") ||
		strings.Contains(out, "0B") ||
		strings.Contains(out, "0 B") ||
		strings.Contains(out, "Total: 0") ||
		strings.Contains(lower, "no cache") ||
		strings.Contains(lower, "empty")
	if !emptyOK {
		t.Fatalf("expected empty-cache indication (0 buckets / 0B / no cache):\n%s", out)
	}
}
```
