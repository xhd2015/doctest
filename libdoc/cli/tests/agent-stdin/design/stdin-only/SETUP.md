# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- Stdin is a pipe with content, no positional args.

## Steps
1. StdinSource = "pipe", Stdin = "design a login form".

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.StdinSource = "pipe"
    req.Stdin = "design a login form"
    return nil
}
```
