# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- Stdin is a pipe with multiline content (simulating heredoc).

## Steps
1. StdinSource = "pipe", Stdin = multiline string with newlines.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.StdinSource = "pipe"
    req.Stdin = "line one\nline two\nline three\n"
    return nil
}
```
