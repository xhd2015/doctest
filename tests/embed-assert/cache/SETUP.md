# Scenario

**Feature**: assert-mod cache materialization is content-addressed and write-once

```
# assert import triggers MaterializeAssertModule
first run -> $CACHE/doctest/assert-mod/<md5>/{assert.go,go.mod}

# no assert import
skip materialization entirely
```

## Preconditions

- Cache root is `$CACHE/doctest/assert-mod/`.
- MD5 matches concatenated `assert/*.go` and `assert/legacy_v1/*.go` sources (sorted, no `*_test.go`).
- Cache leaves run under `lockCacheTests` to avoid cross-test races on `$CACHE`.

## Steps

1. Descendant snapshots cache state, runs doctest, and asserts cache effects.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	lockCacheTests(t)
	req.Env = append(req.Env, "GOWORK=off")
	req.Args = append([]string{"test"}, req.Args...)
	return nil
}
```