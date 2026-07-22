# Scenario

**Feature**: a finished session exists that was created via `CODEX_THREAD_ID` (no `--session-id` flag)

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- A finished session exists that was created via `CODEX_THREAD_ID` (no `--session-id` flag).
- The session has `main_agent_codex_thread_id=codex-fallback-test` but NO `explicit_session_id`.

## Steps
1. Create a session directory with `main_agent_codex_thread_id` set but `explicit_session_id` absent.
2. Write a few events.
3. Run `doctest agent implement --status --session-id codex-fallback-test`.
4. Verify the session is found via codex thread ID fallback.

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

    createdAt := now.Add(-10 * time.Minute)
    meta := map[string]any{
        "main_agent_codex_thread_id": "codex-fallback-test",
        "agent_runner":               "opencode",
        "opencode_session_id":        "ses_fallback_001",
        "created_at":                 createdAt.Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    os.WriteFile(filepath.Join(sessDir, "meta.json"), metaData, 0644)

    nowMs := now.UnixMilli()
    events := []string{
        fmt.Sprintf(`{"type":"step_start","timestamp":%d,"sessionID":"ses_fallback_001","part":{"id":"prt_1","type":"step-start"}}`, nowMs-60000),
        fmt.Sprintf(`{"type":"text","timestamp":%d,"sessionID":"ses_fallback_001","part":{"id":"txt_1","type":"text","text":"work done via codex thread"}}`, nowMs-30000),
    }
    eventsData := strings.Join(events, "\n") + "\n"
    os.WriteFile(filepath.Join(sessDir, "events.jsonl"), []byte(eventsData), 0644)

    req.Args = []string{"agent", "implement", "--status", "--session-id", "codex-fallback-test"}
    return nil
}
```
