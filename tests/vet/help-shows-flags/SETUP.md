# Scenario

**Feature**: the vet command help documents `-v` / multi-dir and `./...` patterns (in-process CLI)

```
cli.RunWithWriter -> doctest vet --help -> usage includes -v, --verbose, <dir...>, ./...
```

## Preconditions

- Vet is registered on the top-level command.
- L2 in-process: `cli.RunWithWriter` captures usage text (no product binary).

## Steps

1. Run `vet --help` via in-process CLI.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"vet", "--help"}
	return nil
}
```
