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
- MD5 matches concatenated `assert/*.go` source (sorted, no `*_test.go`).

## Steps

1. Descendant snapshots cache state, runs doctest, and asserts cache effects.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "GOWORK=off")
	req.Args = append([]string{"test"}, req.Args...)
	return nil
}
```