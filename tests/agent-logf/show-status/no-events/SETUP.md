## Preconditions
- A session directory exists with `meta.json` but no `events.jsonl`.

## Steps
1. Create a session directory with `meta.json` only.
2. Run `doctest agent implement --status --session-id test-status-no-events`.
3. Verify the header block has no timestamps and "No events yet" has a timestamp.

```go
import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    sessionHome := t.TempDir()
    today := time.Now().Format("2006/01/02")
    sessionDir := filepath.Join(sessionHome, today, "sess_test_status_no_ev")
    if err := os.MkdirAll(sessionDir, 0755); err != nil {
        t.Fatalf("create session dir: %v", err)
    }

    meta := map[string]any{
        "explicit_session_id": "test-status-no-events",
        "agent_runner":        "opencode",
        "created_at":           time.Now().Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    if err := os.WriteFile(filepath.Join(sessionDir, "meta.json"), metaData, 0644); err != nil {
        t.Fatalf("write meta.json: %v", err)
    }

    req.Env = append(req.Env, "DOCTEST_DEBUG_SESSION_HOME="+sessionHome)
    req.Args = []string{"agent", "implement", "--status", "--session-id", "test-status-no-events"}
    return nil
}
```
