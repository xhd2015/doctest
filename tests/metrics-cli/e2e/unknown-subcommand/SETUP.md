# Scenario

**Feature**: unknown metrics subcommand exits non-zero (binary)

```
doctest metrics not-a-real-subcmd -> error, exit != 0
```

## Preconditions

- Unknown names are not aliases of valid subcommands.
- L3 e2e: real binary error path.

## Steps

1. Run `metrics` with a bogus subcommand name via product binary.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"metrics", "not-a-real-subcmd"}
	return nil
}
```
