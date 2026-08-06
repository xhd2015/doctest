# Scenario

**Feature**: `-overlay` without a value fails at parse time

```
ParseTestOptions(["-overlay"])
  -> parse error mentioning overlay / argument
```

## Preconditions

- Incomplete flag only; no directory required for arity failure.

## Steps

1. Pass bare `-overlay`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.ParseArgs = []string{"-overlay"}
	return nil
}
```
