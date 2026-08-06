# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- This group tests successful execution of `doctest agent with`.

## Steps
1. Prepend `"agent"`, `"with"`, and `"--agent-runner=opencode"` to the request args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = append([]string{"agent", "with", "--agent-runner=opencode"}, req.Args...)
    return nil
}
```
