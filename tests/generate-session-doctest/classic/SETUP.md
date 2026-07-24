# Scenario

**Feature**: classic `AssembleTestSource` emits the session.Doctest inject contract

```
# classic path
TreeCase -> core.AssembleTestSource -> single-package *_test.go with d inject
```

## Preconditions

- `req.Op` defaults to `"classic"` for all leaves under this group.

## Steps

1. Set `req.Op = "classic"`.
2. Descendants refine AuthorDMode / CasePath / DocTestRoot as needed.

## Context

- Classic inlines Setup/Run/Assert as closures inside `TestXxx`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "classic"
	return nil
}
```
