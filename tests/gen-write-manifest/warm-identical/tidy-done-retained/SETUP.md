# Scenario

**Feature**: warm WriteGoMod retains tidy-done when go.mod/go.sum did not write

```
seed doctest.tidy-done -> second identical WriteGoMod -> tidy-done still present
```

## Steps

1. Inherit warm Setup (includes seeded tidy-done).
2. Confirm tidy-done marker exists before second write.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Mode != "write-gomod-second" {
		t.Fatalf("tidy-done-retained expects Mode write-gomod-second, got %q", req.Mode)
	}
	if !fileExists(filepath.Join(req.GenDir, tidyDoneName)) {
		t.Fatal("parent warm Setup must seed doctest.tidy-done")
	}
	return nil
}
```
