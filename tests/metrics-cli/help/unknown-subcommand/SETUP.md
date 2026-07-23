# Scenario

**Feature**: unknown metrics subcommand exits non-zero (in-process CLI)

```
cli.RunWithWriter -> doctest metrics not-a-real-subcmd -> error, exit != 0
```

## Preconditions

- Unknown names are not aliases of valid subcommands.
- L2 in-process: real CLI error path via `cli.RunWithWriter`.

## Steps

1. Run `metrics` with a bogus subcommand name via in-process CLI.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"metrics", "not-a-real-subcmd"}
	return nil
}
```
