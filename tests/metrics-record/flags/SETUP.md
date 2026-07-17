# Scenario

**Feature**: `doctest test` flag parse for metrics opt-out

```
# CLI args -> core.Options
parseTestOptions([...]) -> Options.NoMetrics

# default metrics on; --no-metrics opts out
```

## Preconditions

- `runner.ParseTestOptions` is exported for package-level tests (or implementer
  re-exports the same parse path used by `runner.Test`).

## Steps

1. Set `req.Op = "parse_flags"` and `req.Args`.
2. Assert `opts.NoMetrics`.

## Context

- Does not run a suite; pure option wiring.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "parse_flags"
	return nil
}
```
