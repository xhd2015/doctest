## Preconditions
- A pre-existing session dir is created with `meta.json` and initial `events.jsonl`.
- `CODEX_THREAD_ID` matches the pre-existing session.
- The mock config produces different events for the resume run.

## Steps
1. Pre-create session dir with `meta.json` (containing `main_agent_codex_thread_id` and `opencode_session_id`).
2. Pre-create `events.jsonl` with an initial event.
3. Set `CODEX_THREAD_ID` to match.
4. Write mock config with a different event for the resume run.
5. Run `doctest agent implement "followup task"`.

```go
import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    threadID := "impl_test_events_append"
    req.Env = append(req.Env, "CODEX_THREAD_ID="+threadID)

    dateDir := time.Now().Format("2006/01/02")
    sessDir := filepath.Join(sessionsDir(), dateDir, "sess_test_events_append")
    if mkErr := os.MkdirAll(sessDir, 0755); mkErr != nil {
        t.Fatalf("create session dir: %v", mkErr)
    }

    meta := map[string]any{
        "main_agent_codex_thread_id": threadID,
        "opencode_session_id":        "saved-sid-events",
        "created_at":                 time.Now().Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    if writeErr := os.WriteFile(filepath.Join(sessDir, "meta.json"), append(metaData, '\n'), 0644); writeErr != nil {
        t.Fatalf("write meta.json: %v", writeErr)
    }

    initialEvent := `{"type":"item.completed","item":{"id":"initial_evt","type":"message","text":"first run output","status":"completed"}}` + "\n"
    if writeErr := os.WriteFile(filepath.Join(sessDir, "events.jsonl"), []byte(initialEvent), 0644); writeErr != nil {
        t.Fatalf("write initial events.jsonl: %v", writeErr)
    }

    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","llm_events":[{"id":"resume_evt","type":"message","text":"resumed and finished"}]}`)
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "followup task"}
    return nil
}
```
