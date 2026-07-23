## Expected

- Subprocess exits 0.
- Stdout/stderr do not indicate session import failure.
- Session-mod cache exists (materialize path used).

## Exit Code

- Exit code 0 after full implementation (RED until then).

```go
import (
	"os"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil && resp.ExitCode == 0 {
		t.Fatalf("unexpected run error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected consumer Once leaf to pass, exit=%d\nstdout:\n%s\nstderr:\n%s\nerr=%v",
			resp.ExitCode, resp.Stdout, resp.Stderr, err)
	}
	cacheDir := expectedSessionCacheDir(t)
	if _, statErr := os.Stat(cacheDir); statErr != nil {
		t.Fatalf("expected session-mod cache after successful session Once leaf: %v", statErr)
	}
}
```
