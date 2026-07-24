# Scenario

**Feature**: `--cold-cache --gen-dir X` allowed when X is outside warm mapping-gen

```
# explicit gen-dir outside warm home
doctest test --cold-cache --gen-dir $TMP/other-cold <tiny-tree>
  -> wipe other-cold on startup
  -> generate into other-cold
  -> leftover remains after finish
```

## Preconditions

- Explicit `--gen-dir` points at a path that is neither equal to nor under
  `$CacheHome/doctest/mapping-gen`.

## Steps

1. Cache sandbox + tiny fixture.
2. Pre-seed marker under the explicit gen dir.
3. Run with `--cold-cache --gen-dir <explicit>`.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	withCacheSandbox(t, req)
	req.CCTestDir = createTempTestProject(t)
	// Explicit dir lives outside the warm home (sibling under sandbox, not under mapping-gen).
	req.CCGenDir = filepath.Join(req.CCCacheHome, "other-cold")
	seedMarker(t, req, req.CCGenDir, "marker-before")
	req.Args = []string{"test", "--cold-cache", "--gen-dir", req.CCGenDir, req.CCTestDir}
	return nil
}
```
