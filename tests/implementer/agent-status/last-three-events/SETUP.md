## Preconditions
- A finished session exists with 8 events in `events.jsonl`.
- `--status` should show all 8 lines in the header but only last 3 in the listing.

## Steps
1. Create a session directory with meta.json.
2. Write 8 events. Events 1-5 should be truncated. Events 6-8 should appear.
3. Run `doctest agent implement --status --session-id status-last-three`.

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
        "explicit_session_id": "status-last-three",
        "agent_runner":        "opencode",
        "created_at":          createdAt.Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    os.WriteFile(filepath.Join(sessDir, "meta.json"), metaData, 0644)

    nowMs := now.UnixMilli()
    events := []string{
        fmt.Sprintf(`{"type":"text","timestamp":%d,"part":{"id":"e1","type":"text","text":"early-event-1"}}`, nowMs-240000),
        fmt.Sprintf(`{"type":"text","timestamp":%d,"part":{"id":"e2","type":"text","text":"early-event-2"}}`, nowMs-210000),
        fmt.Sprintf(`{"type":"text","timestamp":%d,"part":{"id":"e3","type":"text","text":"early-event-3"}}`, nowMs-180000),
        fmt.Sprintf(`{"type":"text","timestamp":%d,"part":{"id":"e4","type":"text","text":"early-event-4"}}`, nowMs-150000),
        fmt.Sprintf(`{"type":"text","timestamp":%d,"part":{"id":"e5","type":"text","text":"early-event-5"}}`, nowMs-120000),
        fmt.Sprintf(`{"type":"tool_use","timestamp":%d,"part":{"type":"tool","tool":"bash","callID":"c1","state":{"status":"completed","input":{"command":"echo six","description":"Sixth event"},"output":"ok","title":"Sixth event"}}}`, nowMs-90000),
        fmt.Sprintf(`{"type":"text","timestamp":%d,"part":{"id":"e7","type":"text","text":"Seventh event text here"}}`, nowMs-60000),
        fmt.Sprintf(`{"type":"text","timestamp":%d,"part":{"id":"e8","type":"text","text":"Eighth event - final output"}}`, nowMs-30000),
    }
    eventsData := strings.Join(events, "\n") + "\n"
    os.WriteFile(filepath.Join(sessDir, "events.jsonl"), []byte(eventsData), 0644)

    req.Args = []string{"agent", "implement", "--status", "--session-id", "status-last-three"}
    return nil
}
```
