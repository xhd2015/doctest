# Scenario

**Feature**: `doctest metrics --help` lists analyze subcommands

```
doctest metrics --help -> usage + path last top summary show prune
```

## Preconditions

- Metrics CLI is registered on the top-level command.

## Steps

1. Run `metrics --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"metrics", "--help"}
	return nil
}
```
