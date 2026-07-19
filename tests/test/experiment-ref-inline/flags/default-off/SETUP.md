# Scenario

**Feature**: experimental ref flag is off by default

```
# parse without --experiment-ref-instead-of-inline
parseTestOptions(["./tests"]) -> ExperimentRefInsteadOfInline=false
```

## Preconditions

- No experiment-ref-related flags in argv.

## Steps

1. Parse a directory-only argv.
2. Expect `ExperimentRefInsteadOfInline` false.

## Context

- Ref generation (P1) enables only when the flag is set; P0 only checks the bool default.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"./tests"}
	return nil
}
```
