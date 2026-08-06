# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- `sh` is available to inspect arguments.

## Steps
1. Run with extra args after the prog, verify they are forwarded.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = append(req.Args, "sh", "-c", "echo $@", "--", "arg1", "arg2")
    return nil
}
```
