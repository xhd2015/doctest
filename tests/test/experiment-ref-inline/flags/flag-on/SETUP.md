# Scenario

**Feature**: `--experiment-ref-instead-of-inline` sets ExperimentRefInsteadOfInline

```
# parse with flag
parseTestOptions(["--experiment-ref-instead-of-inline", "./tests"])
  -> ExperimentRefInsteadOfInline=true
```

## Preconditions

- Flag may appear before the directory operand.

## Steps

1. Parse `--experiment-ref-instead-of-inline` plus a path.
2. Expect field true; remain args still include the path.

## Context

- P0 only sets the bool; ref assembler is P1.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--experiment-ref-instead-of-inline", "./tests"}
	return nil
}
```
