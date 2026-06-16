# Scenario

**Feature**: the main-orchestrator SETUP.md provides shared infrastructure

```
# full TDD cycle: design -> RED -> seal -> implement -> GREEN
orchestrator -> design agent -> writes tests -> RED (all fail)

# seal tests, hand off to implementer
orchestrator -> git add tests/ -> implement agent -> writes code -> GREEN (all pass)

# question/answer loop
user <--questions-- implement agent <--yields-- orchestrator -> resume
```

## Preconditions
- The main-orchestrator SETUP.md provides shared infrastructure.

## Steps
1. Build binaries and provide helpers via parent SETUP.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "TEST_GROUP=full-cycle")
    return nil
}
```
