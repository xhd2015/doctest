# Scenario

**Feature**: parse metrics flags on doctest test

```
parseTestOptions([...]) -> Options.MetricsOn

# default metrics off; --metrics-on opts in
```

## Preconditions

- Uses package parse API (`runner.ParseTestOptions`).

## Steps

1. Set `Op=parse_flags` and `Args`.
2. Assert `opts.MetricsOn`.

## Context

- Flag order before dir is supported.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "parse_flags"
	return nil
}
```
