# Scenario

**Feature**: tidy-done is dropped only when go.mod/go.sum actually wrote

```
seed tidy-done -> go.mod content change -> second WriteGoMod writes go.mod -> tidy-done gone
```

## Steps

1. Inherit content-change Setup (seeds tidy-done, changes parent go.mod).
2. Confirm tidy-done is present before the measured write.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if req.Mode != "write-gomod-second" {
		t.Fatalf("drops-tidy-done expects Mode write-gomod-second, got %q", req.Mode)
	}
	if req.ChangeSourceGoMod == "" {
		t.Fatal("parent must set ChangeSourceGoMod so go.mod actually writes")
	}
	if !fileExists(filepath.Join(req.GenDir, tidyDoneName)) {
		t.Fatal("parent must seed doctest.tidy-done before second WriteGoMod")
	}
	return nil
}
```
