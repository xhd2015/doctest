# Scenario

**Feature**: the agent generate command has scoped help

```
# top-level usage
doctest help -> list subcommands -> stdout

# scoped help
doctest help <subcmd> -> flags, description -> stdout
```

## Preconditions
- The agent generate command has scoped help.

## Steps
1. Run `doctest agent generate --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "generate", "--help"}
    return nil
}
```

