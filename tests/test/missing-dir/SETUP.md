# Scenario

**Feature**: no target directory argument is supplied (L2 CLI)

```
cli.RunWithWriter(["test"]) -> error: test requires <dir>
```

## Preconditions

- Nested L2 root: in-process CLI; no product binary.
- `doctest test` with no directory operand.

## Steps

1. Set `req.Args` to `["test"]` only.
2. Run via `cli.RunWithWriter`; Assert checks message and non-zero exit.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"test"}
	return nil
}
```
