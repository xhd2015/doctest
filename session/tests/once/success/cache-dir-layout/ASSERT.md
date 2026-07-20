## Expected

- Once succeeds.
- `CacheDir` is non-empty and contains path segments `session-once` (under test temp).
- Marker file `probe-write` exists inside CacheDir (writable).
- Returned value is valid JSON.

## Side Effects

- Disk layout includes the once directory used by fn (under t.TempDir).

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
	// Temp layout: .../session-once/<slug>/ (not UserCacheDir/doctest/sessions).
	if !strings.Contains(cd, "session-once") {
		t.Fatalf("cacheDir should contain session-once (t.TempDir layout), got %q", cd)
	}
	base := filepath.Base(cd)
	if base == "" || base == "." || base == "session-once" {
		t.Fatalf("cacheDir should end with slug under session-once, got %q", cd)
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
