## Preconditions
- `DOCTEST_SUBAGENT_AGENT_RUNNER` env var overrides all other detection mechanisms.
- `CODEX_THREAD_ID` is also set to verify priority order.

## Steps
1. Set both `DOCTEST_SUBAGENT_AGENT_RUNNER=idonotexist` and `CODEX_THREAD_ID=abc`.
2. Auto-detection should use env var (priority 1) and pick `"idonotexist"`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Env = append(req.Env,
        "DOCTEST_SUBAGENT_AGENT_RUNNER=idonotexist",
        "CODEX_THREAD_ID=abc",
    )
    return nil
}
```
