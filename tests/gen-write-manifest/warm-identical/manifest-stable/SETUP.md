# Scenario

**Feature**: flush is content-stable when the path→hash map did not change

```
second identical WriteGoMod -> manifest map unchanged -> manifest file not rewritten
```

## Steps

1. Inherit warm Setup (forces old manifest mtime when file exists).
2. Note: under RED (no manifest yet) snapshot may be zero — Assert fails until GREEN.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if req.Mode != "write-gomod-second" {
		t.Fatalf("manifest-stable expects Mode write-gomod-second, got %q", req.Mode)
	}
	// Prefer snapshot when first write already created the manifest.
	// If missing (RED), Assert still requires ManifestExists after Run.
	if fileExists(manifestPath(req.GenDir)) && snapManifestMtimeBefore.IsZero() {
		snapshotManifestMtime(t, req)
	}
	return nil
}
```
