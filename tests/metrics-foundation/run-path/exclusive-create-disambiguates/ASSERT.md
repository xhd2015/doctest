## Expected

- Exactly two paths returned; they are not equal.
- Both basenames match the UTC run pattern prefix `2026-07-17-15-00-00-`.
- Both files exist on disk.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v (resp.Err=%s)", err, resp.Err)
	}
	if len(resp.Paths) != 2 {
		t.Fatalf("got %d paths, want 2: %v", len(resp.Paths), resp.Paths)
	}
	if resp.Paths[0] == resp.Paths[1] {
		t.Fatalf("paths collided: both %q", resp.Paths[0])
	}
	for i, p := range resp.Paths {
		base := filepath.Base(p)
		if !strings.HasPrefix(base, "2026-07-17-15-00-00-") {
			t.Fatalf("path[%d] basename %q missing UTC prefix", i, base)
		}
		if !strings.HasSuffix(base, ".jsonl") {
			t.Fatalf("path[%d] basename %q missing .jsonl", i, base)
		}
		if _, stErr := os.Stat(p); stErr != nil {
			t.Fatalf("path[%d] %q not created: %v", i, p, stErr)
		}
	}
}
```
