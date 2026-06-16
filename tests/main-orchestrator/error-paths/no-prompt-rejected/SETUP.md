# Scenario

**Feature**: `doctest agent implement` requires a prompt

```
# full TDD cycle: design -> RED -> seal -> implement -> GREEN
orchestrator -> design agent -> writes tests -> RED (all fail)

# seal tests, hand off to implementer
orchestrator -> git add tests/ -> implement agent -> writes code -> GREEN (all pass)

# question/answer loop
user <--questions-- implement agent <--yields-- orchestrator -> resume
```

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
