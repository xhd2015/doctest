# Scenario

**Feature**: the mock config returns a non-zero exit code

```
# full TDD cycle: design -> RED -> seal -> implement -> GREEN
orchestrator -> design agent -> writes tests -> RED (all fail)

# seal tests, hand off to implementer
orchestrator -> git add tests/ -> implement agent -> writes code -> GREEN (all pass)

# question/answer loop
user <--questions-- implement agent <--yields-- orchestrator -> resume
```

## Preconditions
- The mock config returns a non-zero exit code.

## Steps
1. Write mock config with exit_code=3 and stderr message.
2. Run `doctest agent implement "broken feature" --agent-runner fake-codex`.

```go
import (
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env, "CODEX_THREAD_ID=impl_test_orch_fail")
    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","session_id":"sess_fail","exit_code":3,"stderr":"build failed: missing dependency"}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "implement broken feature"}
    return nil
}
```
