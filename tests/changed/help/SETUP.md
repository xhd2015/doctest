# Scenario

**Feature**: `--changed` appears in subcommand help output

```
# user reads help for a subcommand
doctest <subcmd> --help -> stdout lists flags including --changed
```

## Preconditions

- The doctest binary is built and available at `req.Bin`.

## Steps

1. Configure `req.Args` to invoke `<subcmd> --help`.
2. Run the doctest binary and capture stdout.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "CHANGED_TEST_GROUP=help")
	return nil
}
```