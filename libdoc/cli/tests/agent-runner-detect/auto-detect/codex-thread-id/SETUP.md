## Preconditions
- `CODEX_THREAD_ID` env var is set but no `DOCTEST_SUBAGENT_AGENT_RUNNER`.
- Auto-detection should match `CODEX_THREAD_ID` (priority 2) and select `"codex"`.

## Steps
1. Set `CODEX_THREAD_ID=abc`.
2. Set `PATH` to a temp directory containing only `cp`, so subagent setup can run but the selected `codex` runner is not found even on machines with Codex installed.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    binDir := t.TempDir()
    if err := os.Symlink("/bin/cp", filepath.Join(binDir, "cp")); err != nil {
        return err
    }
        req.Env = append(req.Env, "DOCTEST_SUBAGENT_AGENT_RUNNER=", "CODEX_THREAD_ID=abc", "PATH="+binDir)
    return nil
}
```
