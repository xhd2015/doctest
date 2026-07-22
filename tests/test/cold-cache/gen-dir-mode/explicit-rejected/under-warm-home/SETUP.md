# Scenario

**Feature**: `--gen-dir` under warm mapping-gen is rejected

```
# subdir of warm home
seed warmHome/foo/marker-before
doctest test --cold-cache --gen-dir $warmHome/foo <tree>
  -> error; marker preserved
```

## Preconditions

- Parent created sandbox + fixture.

## Steps

1. Seed marker under a subdirectory of warm home.
2. Run with `--gen-dir` pointing at that subdirectory.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	under := filepath.Join(req.CCWarmHome, "foo")
	req.CCGenDir = under
	seedMarker(t, req, under, "marker-before")
	req.Args = []string{"test", "--cold-cache", "--gen-dir", under, req.CCTestDir}
	return nil
}
```
