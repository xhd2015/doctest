# Scenario

**Feature**: unified-package assemble (`AssembleUnifiedLeafSource`) uses the same inject contract

```
# unified path
TreeCase -> AssembleUnifiedLeafSource
  -> RunTestLeaf constructs d, no Chdir, no free DOCTEST_* vars
```

## Preconditions

- `req.Op = "unified"`.

## Steps

1. Set Op to unified for all descendants.

## Context

- Unified leaf is a non-test package with `RunTestLeaf` + registry init.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "unified"
	return nil
}
```
