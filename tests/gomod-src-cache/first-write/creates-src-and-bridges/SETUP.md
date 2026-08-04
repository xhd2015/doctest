# Scenario

**Feature**: first cold write creates both source fingerprint and bridges cache files

```
WriteGoModWithVendorBridges (no vendor)
  -> doctest.gomod-src present
  -> doctest.vendor-bridges.json present (empty bridges OK)
```

## Steps

1. Inherit first-write cold fixture.
2. Confirm Mode write-once and empty gen root before measured call.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Mode != "write-once" {
		t.Fatalf("creates-src-and-bridges expects Mode write-once, got %q", req.Mode)
	}
	if req.GenDir == "" || req.ModRoot == "" {
		t.Fatal("parent must prepare GenDir and ModRoot")
	}
	if fileExists(filepath.Join(req.GenDir, gomodSrcName)) {
		t.Fatal("gen root must start without doctest.gomod-src")
	}
	return nil
}
```
