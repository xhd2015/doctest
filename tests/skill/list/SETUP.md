## Preconditions
- No skill name is required for listing.

## Steps
1. Run `doctest skill --list`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"skill", "--list"}
    return nil
}
```

