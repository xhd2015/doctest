## Preconditions
- The doctest command exposes spec documents through the skill subcommand.

## Steps
1. Choose a skill subcommand.
2. Run the doctest command.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 20 * time.Second
    req.Env = append(req.Env, "DOCTEST_SKILL_TEST=1")
    return nil
}
```

