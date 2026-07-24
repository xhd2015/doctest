# Scenario

**Feature**: a session directory exists with `meta.json` and `events.jsonl` containing tool_use events

```
# logf formats agent session events for display
doctest agent logf <session-id> -> reads event file -> formatted text -> stdout

# show-status reports session progress
doctest agent show-status <session-id> -> session state -> stdout
```

## Preconditions
- A session directory exists with `meta.json` and `events.jsonl` containing tool_use events.

## Steps
1. Create a session directory with `meta.json` and `events.jsonl`.
2. Run `doctest agent implement --status --session-id test-status-events`.
3. Verify header block has no timestamps and event display lines have timestamps.

```go
import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "time"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    sessionHome := t.TempDir()
    today := time.Now().Format("2006/01/02")
    sessionDir := filepath.Join(sessionHome, today, "sess_test_status_ev")
    if err := os.MkdirAll(sessionDir, 0755); err != nil {
        t.Fatalf("create session dir: %v", err)
    }

    meta := map[string]any{
        "explicit_session_id": "test-status-events",
        "agent_runner":        "opencode",
        "created_at":           time.Now().Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    if err := os.WriteFile(filepath.Join(sessionDir, "meta.json"), metaData, 0644); err != nil {
        t.Fatalf("write meta.json: %v", err)
    }

    events := []string{
        `{"type":"tool_use","timestamp":1718300000000,"sessionID":"abc","part":{"id":"m1","type":"tool","tool":"bash","callID":"call_1","state":{"status":"completed","title":"run tests","input":{"command":"go test ./..."},"output":"ok"}}}`,
    }
    eventsData := ""
    for _, e := range events {
        eventsData += e + "\n"
    }
    if err := os.WriteFile(filepath.Join(sessionDir, "events.jsonl"), []byte(eventsData), 0644); err != nil {
        t.Fatalf("write events.jsonl: %v", err)
    }

    req.Env = append(req.Env, "DOCTEST_DEBUG_SESSION_HOME="+sessionHome)
    req.Args = []string{"agent", "implement", "--status", "--session-id", "test-status-events"}
    return nil
}
```
