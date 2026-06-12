## Preconditions
- The doc-style test based TDD specification exists under `agents/doctest/doc`.

## Steps
1. Run `doctest skill tdd show`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"skill", "tdd", "show"}
    return nil
}
```
