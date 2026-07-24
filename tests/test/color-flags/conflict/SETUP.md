# Scenario

**Feature**: `--color` and `--no-color` together are rejected (L2)

```
runner.ParseTestOptions([--color, --no-color, .])
  -> error: mutually exclusive
```

## Preconditions

- Both `--color` and `--no-color` are passed (directory unused — parse fails first).

## Steps

1. Set conflicting color flags.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"--color", "--no-color", "."}
	return nil
}
```
