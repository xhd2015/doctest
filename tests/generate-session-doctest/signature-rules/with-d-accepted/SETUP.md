# Scenario

**Feature**: Setup / Run / Assert with `d *session.Doctest` after `t` parse and pass rules

```
# with-d shapes
Setup(t, d *session.Doctest, req *Request) error
Run(t, d *session.Doctest, req *Request) (*Response, error)
Assert(t, d *session.Doctest, req *Request, resp *Response, err error)
  -> Parse*Document OK; rules.Check* OK
```

## Preconditions

- `req.Op = "parse-with-d"`.

## Steps

1. Parse markdown snippets with with-d signatures.
2. Assert ParseErr empty.

## Context

- RED until `libdoc/rules` accept the optional inject param.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Op = "parse-with-d"
	return nil
}
```
