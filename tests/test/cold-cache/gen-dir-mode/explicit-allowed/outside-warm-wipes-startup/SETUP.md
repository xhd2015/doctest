# Scenario

**Feature**: allowed explicit gen-dir is wiped on startup and left after finish

```
# explicit allowed path
pre: other-cold/marker-before
doctest test --cold-cache --gen-dir other-cold <tree>
  -> marker gone; other-cold has gen leftover
```

## Preconditions

- Parent sets `--gen-dir` to a sandbox path outside warm mapping-gen and seeds a marker.

## Steps

1. Inherit Args; assert wipe + leftover on the explicit dir.

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if req.CCGenDir == "" || req.CCMarker == "" {
		t.Fatal("explicit-allowed parent must set GenDir and seed Marker")
	}
	if _, err := os.Stat(req.CCMarker); err != nil {
		t.Fatalf("marker must exist under explicit gen-dir before run: %v", err)
	}
	if len(req.Args) == 0 {
		t.Fatal("req.Args must include --cold-cache --gen-dir")
	}
	return nil
}
```
