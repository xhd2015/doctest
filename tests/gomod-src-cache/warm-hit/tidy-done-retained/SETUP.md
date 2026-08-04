# Scenario

**Feature**: warm hit retains seeded doctest.tidy-done

```
seed tidy-done -> warm second WriteGoModWithVendorBridges -> tidy-done still present
```

## Steps

1. Inherit warm-hit Setup (seeds tidy-done).
2. Confirm tidy-done exists before measured call.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Mode != "write-second" {
		t.Fatalf("tidy-done-retained expects Mode write-second, got %q", req.Mode)
	}
	if !fileExists(filepath.Join(req.GenDir, tidyDoneName)) {
		t.Fatal("parent must seed doctest.tidy-done")
	}
	return nil
}
```
