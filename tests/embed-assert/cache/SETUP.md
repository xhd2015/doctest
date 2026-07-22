# Scenario

**Feature**: assert-mod cache materialization is content-addressed and write-once

```
# always-on for external modules: MaterializeAssertModule before generation
first run -> $DOCTEST_CACHE_HOME|UserCacheDir/doctest/assert-mod/<md5>/{assert.go,go.mod}

# no author assert import still materializes (shared gen-root replace hygiene)
no author import -> still write assert-mod cache entry
```

## Preconditions

- Cache leaves use an isolated `DOCTEST_CACHE_HOME` (temp dir) so they never
  wipe or race the process-global assert-mod cache used by other packages.
- Under that root, layout is `doctest/assert-mod/<md5>/`.
- MD5 matches concatenated `assert/*.go` and `assert/legacy_v1/*.go` sources (sorted, no `*_test.go`).

## Steps

1. Isolate cache home for this leaf and descendants.
2. Descendant snapshots cache state, runs doctest, and asserts cache effects.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func Setup(t *testing.T, req *Request) error {
	isolated := t.TempDir()
	// Child-only isolation via req.Env (subprocess). Never t.Setenv / parent Setenv.
	// Parent-side cache path checks use req.CacheHome, not process env.
	req.CacheHome = isolated
	req.Env = append(req.Env,
		"GOWORK=off",
		core.DoctestCacheHomeEnv+"="+isolated,
	)
	req.Args = append([]string{"test"}, req.Args...)
	return nil
}
```
