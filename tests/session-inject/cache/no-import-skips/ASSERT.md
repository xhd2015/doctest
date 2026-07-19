---
label: heavy
---

## Expected

- Subprocess may pass for the fixture leaf.
- **session-mod** cache for the current key **is created** even though the
  fixture author code does not import session — generated tests always inject
  `d *session.Doctest` and therefore always materialize session-mod.

```go
import (
	"os"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	cacheDir := expectedSessionCacheDir(t)
	if _, statErr := os.Stat(cacheDir); statErr != nil {
		t.Fatalf("session-mod cache should be created for inject even without author session import: %s: %v\nstdout:\n%s\nstderr:\n%s",
			cacheDir, statErr, resp.Stdout, resp.Stderr)
	}
	assertSessionCacheLayout(t, cacheDir)
}
```
