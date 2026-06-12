## Preconditions
- `CODEX_THREAD_ID` is set to simulate a prior session.
- A session dir with matching `meta.json` is pre-created.
- The `meta.json` contains both `codex_thread_id` and `opencode_session_id`.

## Steps
1. Pre-create session dir with `meta.json` containing both thread ID and opencode session ID.
2. Write mock config with completion event.
3. Run `doctest agent implement` with a followup prompt.

```go
import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    threadID := "impl_test_resume_sid"
    req.Env = append(req.Env, "CODEX_THREAD_ID="+threadID)

    home, _ := os.UserHomeDir()
    dateDir := time.Now().Format("2006/01/02")
    sessDir := filepath.Join(home, ".agent-pro", "dedicated-agents", "doctest-agent-implementer", "sessions", dateDir, "sess_test_resume_sid")
    if mkErr := os.MkdirAll(sessDir, 0755); mkErr != nil {
        t.Fatalf("create session dir: %v", mkErr)
    }

    meta := map[string]any{
        "codex_thread_id":     threadID,
        "opencode_session_id": "saved-sid-456",
        "created_at":          time.Now().Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    if writeErr := os.WriteFile(filepath.Join(sessDir, "meta.json"), append(metaData, '\n'), 0644); writeErr != nil {
        t.Fatalf("write meta.json: %v", writeErr)
    }

    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","session_id":"impl_test_resume_sid","stdout_events":[{"type":"item.completed","item":{"id":"m1","type":"message","text":"resumed and finished","status":"completed"}}]}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "followup task"}
    return nil
}
```
