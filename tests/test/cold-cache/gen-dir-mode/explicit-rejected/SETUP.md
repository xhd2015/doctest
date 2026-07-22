# Scenario

**Feature**: `--cold-cache` rejects `--gen-dir` equal to or under warm mapping-gen

```
# reject warm gen-dir
warmHome = $CacheHome/doctest/mapping-gen
doctest test --cold-cache --gen-dir warmHome|warmHome/foo
  -> non-zero exit
  -> stderr Error about default mapping-gen / refuse
  -> do not wipe warm content
```

## Preconditions

- Cache sandbox is active so warm home is under `DOCTEST_CACHE_HOME`.
- Leaves set `--gen-dir` to warm home or a subdirectory and seed a marker that must survive.

## Steps

1. Provide shared sandbox + fixture; leaves specialize gen-dir path.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	withCacheSandbox(t, req)
	req.CCTestDir = createTempTestProject(t)
	return nil
}
```
