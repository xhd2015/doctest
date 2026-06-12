## Preconditions
- `doctest agent implement` requires a prompt.

## Steps
1. Run `doctest agent implement` with no prompt.
2. Expect error.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex"}
    return nil
}
```
