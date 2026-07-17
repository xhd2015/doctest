## Expected

- `RawSourceCacheKeyMD5()` is non-empty.
- Key looks like lowercase hex (length at least 8).
- Embedded content is non-empty.

```go
import (
	"regexp"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("cache key check failed: %v", err)
	}
	if resp.PackageCacheKey == "" {
		t.Fatal("RawSourceCacheKeyMD5 returned empty")
	}
	if !regexp.MustCompile(`^[0-9a-f]{8,}$`).MatchString(resp.PackageCacheKey) {
		t.Fatalf("cache key should be lowercase hex, got %q", resp.PackageCacheKey)
	}
	if resp.ContentLen == 0 {
		t.Fatal("Content() returned empty bytes")
	}
}
```
