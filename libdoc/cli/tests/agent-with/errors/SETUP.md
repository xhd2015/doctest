## Preconditions
- This group tests error conditions for `doctest agent with`.

## Steps
1. Prepend `"agent"` and `"with"` to the request args.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = append([]string{"agent", "with"}, req.Args...)
    return nil
}
```
