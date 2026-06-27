## Expected

- `MaterializeAssertModule` succeeds and returns cache directory path.
- Cache contains `assert.go` and `go.mod` with module path and `go 1.18`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("materialize failed: %v", err)
	}
	if resp.CacheDir == "" {
		t.Fatal("expected non-empty cache dir")
	}
	assertCacheLayoutCore(t, resp.CacheDir)
}
```