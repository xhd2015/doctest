## Preconditions
- `CODEX_THREAD_ID` is set to simulate a prior session.
- A session dir with matching `meta.json` is pre-created.

## Steps
1. Pre-create session dir with `meta.json` containing the thread ID.
2. Write mock config with completion event.
3. Run `doctest agent implement` with answers.

```go
import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    threadID := "impl_test_session_1"
    req.Env = append(req.Env, "CODEX_THREAD_ID="+threadID)

    home, _ := os.UserHomeDir()
    dateDir := time.Now().Format("2006/01/02")
    sessDir := filepath.Join(home, ".agent-pro", "dedicated-agents", "doctest-agent-implementer", "sessions", dateDir, "sess_test_resume")
    if mkErr := os.MkdirAll(sessDir, 0755); mkErr != nil {
        t.Fatalf("create session dir: %v", mkErr)
    }

    meta := map[string]any{
        "codex_thread_id": threadID,
        "created_at":      time.Now().Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    if writeErr := os.WriteFile(filepath.Join(sessDir, "meta.json"), append(metaData, '\n'), 0644); writeErr != nil {
        t.Fatalf("write meta.json: %v", writeErr)
    }

    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","session_id":"impl_test_session_1","stdout_events":[{"type":"item.completed","item":{"id":"m1","type":"message","text":"resumed and done","status":"completed"}}]}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "the port should be 8080"}
    return nil
}
```
