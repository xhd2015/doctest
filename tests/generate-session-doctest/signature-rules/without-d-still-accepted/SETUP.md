# Scenario

**Feature**: legacy without-d Setup/Run/Assert signatures remain accepted

```
# without-d shapes (current)
Setup(t, req *Request) error
Run(t, req *Request) (*Response, error)
Assert(t, req *Request, resp *Response, err error)
  -> still parse OK after P2 (backward compatible)
```

## Preconditions

- `req.Op = "parse-without-d"`.

## Steps

1. Parse without-d signatures.
2. Assert ParseErr empty.

## Context

- Should already be GREEN today; keeps regression coverage after rules widen.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "parse-without-d"
	return nil
}
```
