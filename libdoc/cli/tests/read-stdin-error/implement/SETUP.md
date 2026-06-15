## Preconditions
- This group tests errors from `readStdinIfPresent()` during `doctest agent implement`.

## Steps
1. Prepend `"agent"` and `"implement"` to the request args.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = append([]string{"agent", "implement"}, req.Args...)
	return nil
}
```
