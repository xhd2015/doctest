# Scenario

**Feature**: assert-mod cache materialization is content-addressed and write-once

```
# assert import triggers MaterializeAssertModule
first run -> $DOCTEST_CACHE_HOME|UserCacheDir/doctest/assert-mod/<md5>/{assert.go,go.mod}

# no assert import
skip materialization entirely
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
	t.Setenv(core.DoctestCacheHomeEnv, isolated)
	req.Env = append(req.Env,
		"GOWORK=off",
		core.DoctestCacheHomeEnv+"="+isolated,
	)
	req.Args = append([]string{"test"}, req.Args...)
	return nil
}
```
