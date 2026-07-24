# Scenario

**Feature**: when author declares `d *session.Doctest`, classic keeps that name

```
# author: Setup(t, d *session.Doctest, req *Request)
AssembleTestSource -> generated keeps `d *session.Doctest` (not rewritten to `_`)
```

## Preconditions

- `AuthorDMode=named-d`.

## Steps

1. Assemble classic with named author d.
2. Assert generated params contain `d *session.Doctest` and not forced underscore-only.

## Context

- Name preservation applies to Setup, Run, and Assert closures.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.AuthorDMode = "named-d"
	return nil
}
```
