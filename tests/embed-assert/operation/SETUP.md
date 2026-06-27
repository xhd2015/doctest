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
	req.Env = append(req.Env, "GOWORK=off")
	return nil
}
```