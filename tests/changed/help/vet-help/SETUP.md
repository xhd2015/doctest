# Scenario

**Feature**: `doctest vet --help` documents `--changed`

```
# user reads help for vet subcommand
cli.RunWithWriter -> doctest vet --help -> stdout lists --changed
```

## Preconditions

- In-process CLI via `cli.RunWithWriter` (no product binary).

## Steps

1. Set args to `vet --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"vet", "--help"}
	return nil
}
```
