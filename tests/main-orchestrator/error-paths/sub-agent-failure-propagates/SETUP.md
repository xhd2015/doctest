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
