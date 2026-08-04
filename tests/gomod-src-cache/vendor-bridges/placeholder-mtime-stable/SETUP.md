# Scenario

**Feature**: warm hit does not rewrite overlay placeholder go.mod (mtime stable)

```
force old mtime on placeholder
second WriteGoModWithVendorBridges
  -> placeholder mtime unchanged
```

## Steps

1. Inherit vendor-bridges Setup (forces old placeholder mtime via snapshotBridges).
2. Confirm SnapPlaceholderMtime set.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Mode != "write-second" {
		t.Fatalf("placeholder-mtime-stable expects Mode write-second, got %q", req.Mode)
	}
	if req.SnapPlaceholderPath == "" {
		t.Fatal("parent must snapshot placeholder path")
	}
	if req.SnapPlaceholderMtime.IsZero() {
		t.Fatal("parent must force old placeholder mtime")
	}
	return nil
}
```
