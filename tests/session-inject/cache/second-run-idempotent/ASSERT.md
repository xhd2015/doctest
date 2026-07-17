## Expected

- After second run, `go.mod` bytes equal pre-run snapshot.
- Session-mod layout remains complete.

```go
import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	cacheDir := expectedSessionCacheDir(t)
	after, readErr := os.ReadFile(filepath.Join(cacheDir, "go.mod"))
	if readErr != nil {
		t.Fatalf("read go.mod after second run: %v\nstdout:\n%s\nstderr:\n%s", readErr, resp.Stdout, resp.Stderr)
	}
	if !bytes.Equal(goModBefore, after) {
		t.Fatalf("session-mod go.mod changed on second run:\nbefore:\n%s\nafter:\n%s", goModBefore, after)
	}
	assertSessionCacheLayout(t, cacheDir)
}
```
