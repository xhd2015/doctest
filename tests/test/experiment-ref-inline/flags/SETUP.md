# Scenario

**Feature**: parse `--experiment-ref-instead-of-inline` on doctest test

```
parseTestOptions([...]) -> Options.ExperimentRefInsteadOfInline

# default off; flag opts in
```

## Preconditions

- Uses package parse API (`runner.ParseTestOptions`).

## Steps

1. Set `Op=parse_flags` and `Args`.
2. Assert `opts.ExperimentRefInsteadOfInline`.

## Context

- Flag order before dir is supported (same as other test flags).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "parse_flags"
	return nil
}
```
