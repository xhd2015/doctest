# Scenario

**Feature**: doctest subcommand determines how embedded assert is exercised

```
# doctest test vs doctest build
both must resolve assert via cache replace or -modfile when imported
```

## Preconditions

- Siblings split on doctest operation mode.

## Steps

1. Descendant selects `test` or `build` invocation.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// No process Env (Parallel-safe). Temp fixtures have no go.work.
	_ = req
	return nil
}
```