---
label: heavy
---

## Expected

- `doctest test` exits 0.
- `$CACHE/doctest/assert-mod/<md5>/assert.go` and `go.mod` exist.
- Cached `go.mod` declares `module github.com/xhd2015/doctest/assert` with `go 1.18`.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected first-run materialize test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	cacheDir := expectedAssertCacheDir(t, req.CacheHome)
	assertCacheLayout(t, cacheDir)
}
```