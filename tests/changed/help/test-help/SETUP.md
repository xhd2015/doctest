# Scenario

**Feature**: `doctest test --help` documents `--changed`

```
# user reads help for test subcommand
doctest test --help -> stdout lists --changed
```

## Preconditions

- The doctest binary is built.

## Steps

1. Set args to `test --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"test", "--help"}
	return nil
}
```