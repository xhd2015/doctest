## Preconditions
- A bare `...` pattern is used instead of `./...` or a qualified path.

## Steps
1. Run `doctest vet ...`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"vet", "..."}
	return nil
}
```
