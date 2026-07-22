---
label: heavy
---

## Expected

- Subprocess may pass for the fixture leaf.
- **assert-mod** cache for the current key **is created** even though the
  fixture author code does not import assert — generation always materializes
  assert-mod for external modules (shared gen-root replace hygiene).

```go
import (
	"os"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	cacheDir := expectedAssertCacheDir(t, req.CacheHome)
	if _, statErr := os.Stat(cacheDir); statErr != nil {
		t.Fatalf("assert-mod cache should be created even without author assert import: %s: %v\nstdout:\n%s\nstderr:\n%s",
			cacheDir, statErr, resp.Stdout, resp.Stderr)
	}
	assertCacheLayout(t, cacheDir)
}
```
