## Expected

- Embed script `-cache-key` output matches `assertmod.RawSourceCacheKeyMD5()`.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("cache key check failed: %v", err)
	}
	if resp.ScriptCacheKey == "" {
		t.Fatal("embed script cache key is empty")
	}
	if resp.PackageCacheKey == "" {
		t.Fatal("RawSourceCacheKeyMD5 returned empty string")
	}
	if resp.ScriptCacheKey != resp.PackageCacheKey {
		t.Fatalf("cache key mismatch: script=%s package=%s", resp.ScriptCacheKey, resp.PackageCacheKey)
	}
}
```