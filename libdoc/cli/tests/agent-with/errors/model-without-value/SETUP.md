# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- `--model` flag is present but no value is given.

## Steps
1. Args include `--model` at the end with no following value.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = append(req.Args, "--agent-runner=opencode", "--model")
    return nil
}
```
