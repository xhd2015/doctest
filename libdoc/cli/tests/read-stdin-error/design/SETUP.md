# Scenario

**Feature**: agent design path exercises readStdinIfPresent error handling

```
cli.Run -> agent design -> readStdinIfPresent -> designer.Run
```

## Preconditions
- This group tests errors from `readStdinIfPresent()` during `doctest agent design`.

## Steps
1. Prepend `"agent"` and `"design"` to the request args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = append([]string{"agent", "design"}, req.Args...)
	return nil
}
```
