# Scenario

**Feature**: `doctest test` color flags are parsed at the library boundary (L2)

```
runner.ParseTestOptions([--color, --no-color, .])
  -> mutual exclusion error
```

## Preconditions

- Nested L2 root: `ParseTestOptions` only; no product binary.

## Steps

1. Leaves set `req.Args` for the color-flag scenario under test.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = req
	return nil
}
```
