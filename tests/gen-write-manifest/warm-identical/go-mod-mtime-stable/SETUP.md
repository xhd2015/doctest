# Scenario

**Feature**: warm WriteGoMod does not rewrite go.mod (mtime stable)

```
second WriteGoMod, same content -> go.mod mtime == forced-old before
```

## Steps

1. Inherit warm Setup (first write + old go.mod mtime snapshot).
2. Verify snapshot and Mode before measured second write.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Mode != "write-gomod-second" {
		t.Fatalf("go-mod-mtime-stable expects Mode write-gomod-second, got %q", req.Mode)
	}
	if req.SnapGoModMtimeBefore.IsZero() {
		t.Fatal("parent warm Setup must snapshot go.mod mtime")
	}
	return nil
}
```
