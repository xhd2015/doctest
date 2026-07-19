# Scenario

**Feature**: stderr announces unified (and ref) experimental mode

```
RunTest(..., ExperimentUnifiedPackagePerDoctestTree=true)
  -> stderr/stdout mentions experiment + unified (and ref)
```

## Preconditions

- Same default fixture and unified Options as siblings.

## Steps

1. Run with unified flag on.
2. Expect announce tokens on stderr (stdout fallback).

## Context

- Soft documentation of experimental mode for operators.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "run_gen"
	req.ExperimentUnifiedPackagePerDoctestTree = true
	return nil
}
```
