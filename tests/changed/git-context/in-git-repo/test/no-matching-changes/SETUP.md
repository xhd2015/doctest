# Scenario

**Feature**: no doctest changes prints warning and exits 0

```
# no affected doctest files
doctest test --changed -> silent stderr -> exit 0 | with -v -> stderr shows "--changed: 0 tests"
```

## Preconditions

- Git repo with committed fixture tree.

## Steps

1. Prepare baseline repo (descendant leaf applies or omits extra changes).
2. Run `doctest test --changed`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "CHANGED_SCENARIO=no-matching")
	return nil
}
```