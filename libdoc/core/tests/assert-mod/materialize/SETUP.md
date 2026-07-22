# Scenario

**Feature**: MaterializeAssertModule writes content-addressed assert cache

```
# write-once cache
MaterializeAssertModule -> $CACHE/doctest/assert-mod/<md5>/{assert.go,go.mod}
```

## Preconditions

- Siblings test first-write layout and second-call idempotency.

## Steps

1. Descendant sets `runKind = "materialize"`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.RunKind = "materialize"
	req.ModPath = "example.com/app"
	return nil
}
```