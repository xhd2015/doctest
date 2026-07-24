# Scenario

**Feature**: classic assemble constructs `d *session.Doctest`, passes it, and drops Chdir/free vars

```
# author declares d (named-d)
AssembleTestSource
  -> d := &session.Doctest{ROOT, CASE, SESSION_ID}
  -> setup(t, d, req); run(t, d, req); assert(t, d, req, resp, err)
  -> no os.Chdir; no package free DOCTEST_* vars
```

## Preconditions

- Author signatures use the required with-d shape (`AuthorDMode=named-d`).
- Leaf path is the default `nested/leaf`.

## Steps

1. Assemble classic with default author params.
2. Assert inject contract on generated source.

## Context

- This is the primary RED leaf against current Chdir + free-var boilerplate.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.AuthorDMode = "named-d"
	return nil
}
```
