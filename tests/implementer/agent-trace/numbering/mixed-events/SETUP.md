# Scenario

**Feature**: a finished session exists with 5 events in `events.jsonl` where 3 are displayable and 2 are non-displayable system events (step_start)

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- A finished session exists with 5 events in `events.jsonl` where 3 are displayable and 2 are non-displayable system events (step_start).
- The 2 system events should NOT consume trace numbers; only the 3 visible events should be numbered [1],[2],[3].

## Steps
1. Create a session directory with meta.json.
2. Write 5 events: 2 step_start (non-displayable) interleaved with 3 text/tool_use (displayable).
3. Run `doctest agent implement --session-id trace-mixed-1 --trace`.

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
    sessHome := sessionsDir(req)
    now := time.Now()
    dateDir := now.Format("2006/01/02")
    sessName := fmt.Sprintf("sess_%s_%d", now.Format("150405"), now.UnixNano())
    sessDir := filepath.Join(sessHome, dateDir, sessName)
    os.MkdirAll(sessDir, 0755)

    meta := map[string]any{
        "explicit_session_id": "trace-mixed-1",
        "created_at":          now.Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    os.WriteFile(filepath.Join(sessDir, "meta.json"), metaData, 0644)

    nowMs := now.UnixMilli()
    events := []string{
        fmt.Sprintf(`{"type":"step_start","timestamp":%d,"part":{"id":"s1","type":"step-start"}}`, nowMs-240000),
        fmt.Sprintf(`{"type":"text","timestamp":%d,"part":{"id":"e1","type":"message","text":"First displayable message"}}`, nowMs-200000),
        fmt.Sprintf(`{"type":"step_start","timestamp":%d,"part":{"id":"s2","type":"step-start"}}`, nowMs-180000),
        fmt.Sprintf(`{"type":"tool_use","timestamp":%d,"part":{"type":"tool","tool":"bash","callID":"c1","state":{"status":"completed","input":{"description":"Build project"},"output":"ok","title":"Build project"}}}`, nowMs-120000),
        fmt.Sprintf(`{"type":"text","timestamp":%d,"part":{"id":"e2","type":"message","text":"Final output message"}}`, nowMs-60000),
    }
    eventsData := strings.Join(events, "\n") + "\n"
    os.WriteFile(filepath.Join(sessDir, "events.jsonl"), []byte(eventsData), 0644)

    req.Args = []string{"agent", "implement", "--session-id", "trace-mixed-1", "--trace"}
    return nil
}
```
