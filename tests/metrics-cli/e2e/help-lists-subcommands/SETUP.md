# Scenario

**Feature**: `doctest metrics --help` lists analyze subcommands (binary)

```
doctest metrics --help -> usage + path last top summary show prune
```

## Preconditions

- Metrics CLI is registered on the top-level command.
- L3 e2e: real binary argv/help text.

## Steps

1. Run `metrics --help` via product binary.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"metrics", "--help"}
	return nil
}
```
