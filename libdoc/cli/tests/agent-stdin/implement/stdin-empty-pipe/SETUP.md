# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- Stdin is a pipe but with empty content (EOF immediately).

## Steps
1. StdinSource = "pipe", Stdin = "" (empty string, pipe closed with no data).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.StdinSource = "pipe"
    req.Stdin = ""
    return nil
}
```
