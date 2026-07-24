# Scenario

**Feature**: when author omits `d *session.Doctest`, assemble fails clearly (no auto-inject)

```
# author: Setup(t, req) without d
AssembleTestSource -> error: missing d *session.Doctest (no auto-inject)
```

## Preconditions

- `AuthorDMode=omit` (explicit).

## Steps

1. Assemble classic with omitted author d.
2. Assert assemble returns a clear missing-d error (not silent `_` inject).

## Context

- Broken author signatures must be fixed in the test tree, not papered over by gen.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.AuthorDMode = "omit"
	return nil
}
```
