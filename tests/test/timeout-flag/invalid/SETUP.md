# Scenario

**Feature**: invalid `-timeout` values are rejected at parse time (L2)

```
runner.ParseTestOptions([-timeout, bogus, .]) -> error mentioning timeout
```

## Preconditions

- Nested L2 root: parse only; directory unused after parse fails.

## Steps

1. Set args with invalid duration.
2. Assert non-zero exit mapping and timeout in error text.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"-timeout", "bogus", "."}
	return nil
}
```
