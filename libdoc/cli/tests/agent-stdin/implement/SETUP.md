# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- This group tests `doctest agent implement` command.

## Steps
1. Prepend `"agent"` and `"implement"` to the request args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = append([]string{"agent", "implement"}, req.Args...)
    return nil
}
```
