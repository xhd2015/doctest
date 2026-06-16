# Scenario

**Feature**: a finished session directory exists with full `meta.json` and 5 events in `events.jsonl`

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- A finished session directory exists with full `meta.json` and 5 events in `events.jsonl`.
- Events use the actual opencode JSON format with `part` field and `timestamp` in milliseconds.
- Only the last 3 events should appear in the output listing.

## Steps
1. Create a session directory with all metadata fields.
2. Write 5 events. Event 2 content should be truncated (only last 3 shown). Events 3-5 should appear.
3. Run `doctest agent implement --status --session-id status-finished-test`.

```go
import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    sessHome := sessionsDir()
    now := time.Now()
    dateDir := now.Format("2006/01/02")
    sessName := fmt.Sprintf("sess_%s_%d", now.Format("150405"), now.UnixNano())
    sessDir := filepath.Join(sessHome, dateDir, sessName)
    os.MkdirAll(sessDir, 0755)

    createdAt := now.Add(-10 * time.Minute)
    meta := map[string]any{
        "explicit_session_id":                "status-finished-test",
        "agent_runner":                       "opencode",
        "main_agent_codex_thread_id":         "codex-thread-abc123",
        "opencode_session_id":                "ses_open_xyz789",
        "created_at":                         createdAt.Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    os.WriteFile(filepath.Join(sessDir, "meta.json"), metaData, 0644)

    nowMs := now.UnixMilli()
    events := []string{
        fmt.Sprintf(`{"type":"step_start","timestamp":%d,"sessionID":"ses_open_xyz789","part":{"id":"prt_1","type":"step-start"}}`, nowMs-120000),
        fmt.Sprintf(`{"type":"text","timestamp":%d,"sessionID":"ses_open_xyz789","part":{"id":"txt_trunc","type":"text","text":"THIS-SHOULD-BE-TRUNCATED"}}`, nowMs-100000),
        fmt.Sprintf(`{"type":"tool_use","timestamp":%d,"sessionID":"ses_open_xyz789","part":{"type":"tool","tool":"bash","callID":"call_2","state":{"status":"completed","input":{"command":"go build ./...","description":"Build project"},"output":"ok","title":"Build project"}}}`, nowMs-90000),
        fmt.Sprintf(`{"type":"tool_use","timestamp":%d,"sessionID":"ses_open_xyz789","part":{"type":"tool","tool":"task","callID":"call_3","state":{"status":"completed","input":{"description":"Run tests","prompt":"Run go test ./... and report results"},"output":"All 16 tests pass","title":"Run tests"}}}`, nowMs-30000),
        fmt.Sprintf(`{"type":"text","timestamp":%d,"sessionID":"ses_open_xyz789","part":{"id":"txt_2","type":"text","text":"All 16 tests pass — implementation complete."}}`, nowMs-10000),
    }
    eventsData := strings.Join(events, "\n") + "\n"
    os.WriteFile(filepath.Join(sessDir, "events.jsonl"), []byte(eventsData), 0644)

    req.Args = []string{"agent", "implement", "--status", "--session-id", "status-finished-test"}
    return nil
}
```
