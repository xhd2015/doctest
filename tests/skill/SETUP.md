# Scenario

**Feature**: the doctest command exposes spec documents through the skill subcommand

```
# expose embedded spec documents
doctest skill <name> --show -> embedded .md doc -> stdout

# list available skills
doctest skill --list -> skill names -> stdout
```

## Preconditions
- The doctest command exposes spec documents through the skill subcommand.

## Steps
1. Choose a skill subcommand.
2. Run the doctest command.

```go
import (
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 20 * time.Second
    req.Env = append(req.Env, "DOCTEST_SKILL_TEST=1")
    return nil
}
```

