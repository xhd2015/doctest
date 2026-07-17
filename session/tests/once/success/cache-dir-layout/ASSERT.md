## Expected

- Once succeeds.
- `CacheDir` is non-empty and contains path segments `doctest/sessions` and `once-`.
- `CacheDir` is under the user cache sessions root.
- Marker file `probe-write` exists inside CacheDir (writable).
- Returned value is valid JSON.

## Side Effects

- Disk layout includes the once directory used by fn.

```go
import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("Once error: %v", resp.Err)
	}
	if resp.CacheDir == "" {
		t.Fatal("expected non-empty cacheDir passed to fn")
	}
	cd := filepath.Clean(resp.CacheDir)
	if !strings.Contains(cd, filepath.Join("doctest", "sessions")) {
		t.Fatalf("cacheDir should be under doctest/sessions, got %q", cd)
	}
	base := filepath.Base(cd)
	if !strings.HasPrefix(base, "once-") {
		t.Fatalf("cacheDir base should be once-<slug>, got %q", base)
	}
	root := userCacheSessionsRoot(t)
	if !strings.HasPrefix(cd, filepath.Clean(root)) {
		t.Fatalf("cacheDir %q not under sessions root %q", cd, root)
	}
	marker := filepath.Join(cd, "probe-write")
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("cacheDir not writable or marker missing: %v", err)
	}
	if !json.Valid(resp.Value) {
		t.Fatalf("value not valid JSON: %s", resp.Value)
	}
	if resp.FnCalls != 1 {
		t.Fatalf("fn calls=%d want 1", resp.FnCalls)
	}
}
```
