# Scenario

**Feature**: unified leaf injects `d *session.Doctest` without Chdir or free vars

```
# AssembleUnifiedLeafSource
func RunTestLeaf(t *testing.T) {
  d := &session.Doctest{...}
  setup / Run / assert receive d
  no os.Chdir; no package free DOCTEST_* vars
}
```

## Preconditions

- Author omits d; leaf path default.

## Steps

1. Assemble unified leaf source.
2. Assert inject contract.

## Context

- Unified currently mirrors ref leaf Chdir + free vars → RED until P2.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.AuthorDMode = "named-d"
	return nil
}
```
