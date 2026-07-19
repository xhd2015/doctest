# Scenario

**Feature**: auto cold home is wiped on startup, kept after finish, and announced

```
# startup wipe + leftover + announce
pre: coldHome/marker-before exists
doctest test --cold-cache <tree>
  -> marker-before gone
  -> coldHome still exists with generated content
  -> stderr mentions cold-cache (gen path / GOCACHE / count)
```

## Preconditions

- Parent auto setup seeds the marker and configures `--cold-cache` without `--gen-dir`.

## Steps

1. Inherit parent Args/Env; assert wipe, leftover, and announcement.

```go
import (
	"os"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	// Parent configured --cold-cache (auto gen) and seeded the marker.
	// Verify preconditions so Assert can rely on wipe semantics.
	if st.ColdHome == "" || st.Marker == "" {
		t.Fatal("auto parent must set ColdHome and seed Marker")
	}
	if _, err := os.Stat(st.Marker); err != nil {
		t.Fatalf("marker must exist before cold-cache run: %v", err)
	}
	if len(req.Args) == 0 {
		t.Fatal("req.Args must be set by auto parent")
	}
	return nil
}
```
