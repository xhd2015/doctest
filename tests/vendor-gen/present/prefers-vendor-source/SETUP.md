# Scenario

**Feature**: gen resolution prefers vendored sources via replace to vendor path

```
vendor/example.com/dep/dep.go contains DistinctiveMarker
  -> WriteGoMod
  -> replace example.com/dep => <modRoot>/vendor/example.com/dep
  -> reading replace target surfaces DistinctiveMarker (not module cache)
```

## Steps

1. Inherit present fixture with distinctive marker in dep package.
2. After WriteGoMod, assert replace line and that marker is readable at the
   replace target path (unit-level preference proof without full `go test` e2e).

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.DistinctiveMarker == "" {
		t.Fatal("DistinctiveMarker must be set by seedTinyVendor")
	}
	src := filepath.Join(vendorModuleDir(req, req.SampleModPath), "dep.go")
	if !fileExists(src) {
		t.Fatalf("missing vendored source %s", src)
	}
	if !strings.Contains(readFileOrEmpty(src), req.DistinctiveMarker) {
		t.Fatalf("vendored source missing marker %q", req.DistinctiveMarker)
	}
	return nil
}
```
