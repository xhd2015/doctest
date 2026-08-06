---
label: e2e
---

## Expected

- Subprocess `doctest test` exits 0 (or at least reaches materialize — prefer exit 0 once full stack works).
- `$UserCacheDir/doctest/session-mod/<md5>/go.mod` exists with
  `module github.com/xhd2015/doctest/session`.
- At least one `.go` source file exists in the cache dir.

## Exit Code

- Prefer exit code 0 after implementation; RED may fail compile until session package exists in embed.

```go
import (
	"os"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	// Materialize should occur even if leaf later fails; require cache layout.
	cacheDir := expectedSessionCacheDir(t)
	if _, statErr := os.Stat(cacheDir); statErr != nil {
		t.Fatalf("expected session-mod cache after run with session import: %v\nstdout:\n%s\nstderr:\n%s\nrunErr=%v exit=%d",
			statErr, resp.Stdout, resp.Stderr, err, resp.ExitCode)
	}
	assertSessionCacheLayout(t, cacheDir)
}
```
