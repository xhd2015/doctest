# Scenario

**Feature**: `doctest metrics --help` lists analyze subcommands (in-process CLI)

```
cli.RunWithWriter -> doctest metrics --help -> usage + path last top summary show prune
```

## Preconditions

- Metrics CLI is registered on the top-level command.
- L2 in-process: `cli.RunWithWriter` captures usage text (no product binary).

## Steps

1. Run `metrics --help` via in-process CLI.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"metrics", "--help"}
	return nil
}
```
