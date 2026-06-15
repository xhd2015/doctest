## Preconditions
- This group tests `doctest agent design` command.

## Steps
1. Prepend `"agent"` and `"design"` to the request args.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = append([]string{"agent", "design"}, req.Args...)
    return nil
}
```
