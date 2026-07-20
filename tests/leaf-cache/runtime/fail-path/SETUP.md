# Scenario

**Feature**: failed leaves are never stored as passes

```
# fail fixture
doctest test fail-tree -> exit != 0, no PutPass
doctest test fail-tree -> still executes, 0 Cached
```

## Preconditions

- Fixture has at least one always-fail leaf.
- Leaf-cache enabled (no disable flags).

## Steps

1. Prepare fail fixture.
2. Run twice with default args; both fail with 0 Cached.

## Context

- Complements warm pass-path; ensures PutPass is success-only.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	t.Helper()
	req.FixtureDir = prepareFailFixture(t, 1)
	req.Args = []string{"test", req.FixtureDir}
	req.Args2 = []string{"test", req.FixtureDir}
	return nil
}
```
