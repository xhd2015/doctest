# Scenario

**Feature**: no previous command arguments are required

```
# top-level usage
doctest help -> list subcommands -> stdout

# scoped help
doctest help <subcmd> -> flags, description -> stdout
```

## Preconditions
- No previous command arguments are required.

## Steps
1. Run `doctest --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"--help"}
    return nil
}
```

