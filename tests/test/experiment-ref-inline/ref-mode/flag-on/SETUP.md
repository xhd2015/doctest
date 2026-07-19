# Scenario

**Feature**: flag on uses ref package DAG (root package + thin leaves)

```
RunTest(2-leaf, ExperimentRefInsteadOfInline=true, GenDir=tmp)
  -> ref assemble path
  -> root owns Run/types/helpers once; leaves import
```

## Preconditions

- Field explicitly true.
- P1 implementer wires generation when the option is set.

## Steps

1. Set `ExperimentRefInsteadOfInline=true`.
2. Leaf asserts pass and/or layout.

## Context

- Sibling leaves under this node are MECE by assertion focus (pass vs marker once vs thin import vs stderr).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.ExperimentRefInsteadOfInline = true
	return nil
}
```
