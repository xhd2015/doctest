# Scenario

**Feature**: `doctest vet --help` documents `--changed`

```
# user reads help for vet subcommand
doctest vet --help -> stdout lists --changed
```

## Preconditions

- The doctest binary is built.

## Steps

1. Set args to `vet --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"vet", "--help"}
	return nil
}
```