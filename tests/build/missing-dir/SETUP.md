## Preconditions
- No target directory argument is supplied.

## Steps
1. Run `doctest build`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"build"}
    return nil
}
```

