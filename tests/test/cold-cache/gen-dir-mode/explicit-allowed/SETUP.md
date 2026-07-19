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

func Setup(t *testing.T, req *Request) error {
	withCacheSandbox(t, req)
	st.TestDir = createTempTestProject(t)
	// Explicit dir lives outside the warm home (sibling under sandbox, not under mapping-gen).
	st.GenDir = filepath.Join(st.CacheHome, "other-cold")
	seedMarker(t, st.GenDir, "marker-before")
	req.Args = []string{"test", "--cold-cache", "--gen-dir", st.GenDir, st.TestDir}
	return nil
}
```
