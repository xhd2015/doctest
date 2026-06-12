## Preconditions
- Build has scoped help with runner options.

## Steps
1. Run `doctest build --help`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"build", "--help"}
    return nil
}
```
