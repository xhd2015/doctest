# Scenario

**Feature**: three session directories exist with different `explicit_session_id`, `agent_runner`, and `created_at`

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- Three session directories exist with different `explicit_session_id`, `agent_runner`, and `created_at`.
- The `--list-sessions` flag lists all sessions found within the last 7 days.

## Steps
1. Create 3 session directories with distinct metadata.
2. Run `doctest agent implement --list-sessions`.

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

    sessions := []struct {
        id       string
        runner   string
        codexID  string
        openID   string
        createdAt time.Time
    }{
        {"list-alpha", "opencode", "codex-a", "ses-a", now.Add(-30 * time.Minute)},
        {"list-beta",  "codex",    "codex-b", "ses-b", now.Add(-2 * time.Hour)},
        {"list-gamma", "opencode", "codex-c", "ses-c", now.Add(-1 * time.Hour)},
    }

    for _, s := range sessions {
        dateDir := s.createdAt.Format("2006/01/02")
        sessName := fmt.Sprintf("sess_%s_%d", s.createdAt.Format("150405"), s.createdAt.UnixNano())
        sessDir := filepath.Join(sessHome, dateDir, sessName)
        os.MkdirAll(sessDir, 0755)

        meta := map[string]any{
            "explicit_session_id":                s.id,
            "agent_runner":                       s.runner,
            "main_agent_codex_thread_id":         s.codexID,
            "opencode_session_id":                s.openID,
            "created_at":                         s.createdAt.Format(time.RFC3339),
        }
        metaData, _ := json.MarshalIndent(meta, "", "  ")
        os.WriteFile(filepath.Join(sessDir, "meta.json"), metaData, 0644)
    }

    req.Args = []string{"agent", "implement", "--list-sessions"}
    return nil
}
```
