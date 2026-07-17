---
label: heavy
---

## Expected

- Subprocess may pass or fail for unrelated reasons; **session-mod** cache for
  the current key must **not** be created solely by this run.
- After run, expected session cache dir does not exist.

```go
import (
	"os"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	cacheDir := expectedSessionCacheDir(t)
	if _, statErr := os.Stat(cacheDir); statErr == nil {
		t.Fatalf("session-mod cache should not be created without session import: %s\nstdout:\n%s\nstderr:\n%s",
			cacheDir, resp.Stdout, resp.Stderr)
	}
}
```
