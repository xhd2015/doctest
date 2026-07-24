# Scenario

**Feature**: `doctest test --help` documents `--changed`

```
# user reads help for test subcommand
cli.RunWithWriter -> doctest test --help -> stdout lists --changed
```

## Preconditions

- In-process CLI via `cli.RunWithWriter` (no product binary).

## Steps

1. Set args to `test --help`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"test", "--help"}
	return nil
}
```
