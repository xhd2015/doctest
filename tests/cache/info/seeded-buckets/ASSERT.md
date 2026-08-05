## Expected

- Exit code 0.
- Stdout includes injectable Cache home and Doctest root paths.
- Each seeded bucket name appears (`leaf-cache`, `mapping-gen`).
- A Total line is present with a non-zero human size (K/M/G/B unit).
- Output uses human size units (not raw byte-only dump without units).

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
	mustContain(t, out, req.CacheHome, "Cache home path")
	mustContain(t, out, req.DoctestRoot, "Doctest root path")
	for _, b := range req.SeededBuckets {
		mustContain(t, out, b, "bucket name")
	}
	mustContain(t, out, "Total", "total line")
	if !hasHumanSizeUnit(out) {
		t.Fatalf("expected human size units (B/K/M/G) in stdout:\n%s", out)
	}
	// Total should be non-zero for seeded data (~5K+).
	lower := strings.ToLower(out)
	if strings.Contains(out, "Total: 0") || strings.Contains(out, "Total:0") ||
		strings.Contains(lower, "total: 0b") || strings.Contains(lower, "total: 0 b") {
		t.Fatalf("expected non-zero Total for seeded buckets:\n%s", out)
	}
}
```
