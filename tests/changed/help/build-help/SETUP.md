# Scenario

**Feature**: `doctest build --help` documents `--changed`

```
# user reads help for build subcommand
doctest build --help -> stdout lists --changed
```

## Preconditions

- The doctest binary is built.

## Steps

1. Set args to `build --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"build", "--help"}
	return nil
}
```