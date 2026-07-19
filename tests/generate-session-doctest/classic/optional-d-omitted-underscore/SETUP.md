# Scenario

**Feature**: when author omits the inject param, classic inserts `_ *session.Doctest`

```
# author: Setup(t, req) without d
AssembleTestSource -> generated closures include `_ *session.Doctest` as second param
# call sites still pass the real d value
```

## Preconditions

- `AuthorDMode=omit` (explicit).

## Steps

1. Assemble classic with omitted author d.
2. Assert `_ *session.Doctest` appears in generated source and call sites pass `d`.

## Context

- Distinct from injects-d: focuses on the underscore placeholder in signatures.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.AuthorDMode = "omit"
	return nil
}
```
