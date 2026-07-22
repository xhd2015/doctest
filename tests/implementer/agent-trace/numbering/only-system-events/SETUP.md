# Scenario

**Feature**: a finished session exists with only non-displayable system events (step_start, step_finish) in `events.jsonl`

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- A finished session exists with only non-displayable system events (step_start, step_finish) in `events.jsonl`.
- Since no event produces visible output, no trace numbers at all should appear.

## Steps
1. Create a session directory with meta.json.
2. Write 3 events all of type step_start (non-displayable).
3. Run `doctest agent implement --session-id trace-system-only --trace`.

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
        "explicit_session_id": "trace-system-only",
        "created_at":          now.Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    os.WriteFile(filepath.Join(sessDir, "meta.json"), metaData, 0644)

    nowMs := now.UnixMilli()
    events := []string{
        fmt.Sprintf(`{"type":"step_start","timestamp":%d,"part":{"id":"s1","type":"step-start"}}`, nowMs-120000),
        fmt.Sprintf(`{"type":"step_start","timestamp":%d,"part":{"id":"s2","type":"step-start"}}`, nowMs-60000),
        fmt.Sprintf(`{"type":"step_start","timestamp":%d,"part":{"id":"s3","type":"step-start"}}`, nowMs-30000),
    }
    eventsData := strings.Join(events, "\n") + "\n"
    os.WriteFile(filepath.Join(sessDir, "events.jsonl"), []byte(eventsData), 0644)

    req.Args = []string{"agent", "implement", "--session-id", "trace-system-only", "--trace"}
    return nil
}
```
