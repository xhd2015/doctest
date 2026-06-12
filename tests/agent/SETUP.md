## Preconditions
- Agent commands must be able to use fake Codex instead of a real LLM.

## Steps
1. Lookup `fake-codex` from PATH; skip if not installed.
2. Write a deterministic fake Codex script.
3. Configure the doctest command to use fake Codex.

```go
import (
    "os"
    "os/exec"
    "path/filepath"
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    tmp := t.TempDir()
    fakeCodex, err := exec.LookPath("fake-codex")
    if err != nil {
        t.Skip("fake-codex not in PATH; install via: go install github.com/xhd2015/agent-pro/cmd/fake-codex@latest")
        return nil
    }
    script := filepath.Join(tmp, "fake-codex-script.json")
    content := `{"events":[{"type":"item.completed","item":{"id":"m1","type":"message","text":"fake doctest agent completed","status":"completed"}}]}`
    if err := os.WriteFile(script, []byte(content), 0644); err != nil {
        t.Fatalf("write fake script: %v", err)
    }
    req.Env = append(req.Env,
        "AGENT_RUNNER_FAKE_CODEX_PATH="+fakeCodex,
        "FAKE_CODEX_SCRIPT="+script,
    )
    req.Timeout = 60 * time.Second
    return nil
}
```
