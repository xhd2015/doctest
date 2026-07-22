# Scenario

**Feature**: a session directory exists with `meta.json` but an empty `events.jsonl`

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- A session directory exists with `meta.json` but an empty `events.jsonl`.
- The `--status` flag should report "No events yet".

## Steps
1. Create a session directory with meta.json.
2. Write an empty events.jsonl (or with only whitespace).
3. Run `doctest agent implement --status --session-id status-no-events-test`.

```go
import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
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
        "explicit_session_id": "status-no-events-test",
        "created_at":          now.Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    os.WriteFile(filepath.Join(sessDir, "meta.json"), metaData, 0644)

    os.WriteFile(filepath.Join(sessDir, "events.jsonl"), []byte(""), 0644)

    req.Args = []string{"agent", "implement", "--status", "--session-id", "status-no-events-test"}
    return nil
}
```
