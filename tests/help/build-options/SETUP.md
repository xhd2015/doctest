# Scenario

**Feature**: build has scoped help with runner options

```
# top-level usage
doctest help -> list subcommands -> stdout

# scoped help
doctest help <subcmd> -> flags, description -> stdout
```

## Preconditions
- Build has scoped help with runner options.

## Steps
1. Run `doctest build --help`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = []string{"build", "--help"}
    return nil
}
```
