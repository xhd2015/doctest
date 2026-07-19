# Scenario

**Feature**: parse `--experiment-unified-package-per-doctest-tree` on doctest test

```
parseTestOptions([...]) -> Options.ExperimentUnifiedPackagePerDoctestTree
                        -> Options.ExperimentRefInsteadOfInline when unified on

# default off; flag opts in and forces ref
```

## Preconditions

- Uses package parse API (`runner.ParseTestOptions`).

## Steps

1. Set `Op=parse_flags` and `Args`.
2. Assert unified (and ref when flag on).

## Context

- Flag order before dir is supported (same as other test flags).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "parse_flags"
	return nil
}
```
