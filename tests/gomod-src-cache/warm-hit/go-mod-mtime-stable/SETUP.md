# Scenario

**Feature**: warm hit does not rewrite gen go.mod (mtime stable)

```
second WriteGoModWithVendorBridges, same sources
  -> go.mod mtime == forced-old before
```

## Steps

1. Inherit warm-hit Setup (first write + old go.mod mtime).
2. Verify snapshot present before measured second call.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Mode != "write-second" {
		t.Fatalf("go-mod-mtime-stable expects Mode write-second, got %q", req.Mode)
	}
	if req.SnapGoModMtimeBefore.IsZero() {
		t.Fatal("parent warm Setup must snapshot go.mod mtime")
	}
	return nil
}
```
