## Preconditions
- Test has scoped help with runner options.

## Steps
1. Run `doctest test --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"test", "--help"}
    return nil
}
```
