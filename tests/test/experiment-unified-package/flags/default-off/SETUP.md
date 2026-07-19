# Scenario

**Feature**: experimental unified-package flag is off by default

```
# parse without --experiment-unified-package-per-doctest-tree
parseTestOptions(["./tests"])
  -> ExperimentUnifiedPackagePerDoctestTree=false
  -> ExperimentRefInsteadOfInline=false
```

## Preconditions

- No experiment-unified or experiment-ref flags in argv.

## Steps

1. Parse a directory-only argv.
2. Expect both experiment fields false.

## Context

- Unified generation enables only when the flag is set.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"./tests"}
	return nil
}
```
