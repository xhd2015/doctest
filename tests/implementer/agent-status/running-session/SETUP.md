# Scenario

**Feature**: a session directory exists with `meta.json` and a `pid` file pointing to the current process

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- A session directory exists with `meta.json` and a `pid` file pointing to the current process.
- The `--status` flag should show "running" status with the PID.

## Steps
1. Create a session directory with meta.json.
2. Write a pid file containing the current process PID.
3. Run `doctest agent implement --status --session-id status-running-test`.

```go
import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
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
        "explicit_session_id": "status-running-test",
        "created_at":          now.Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    os.WriteFile(filepath.Join(sessDir, "meta.json"), metaData, 0644)

    pidStr := strconv.Itoa(os.Getpid())
    os.WriteFile(filepath.Join(sessDir, "pid"), []byte(pidStr), 0644)

    req.Args = []string{"agent", "implement", "--status", "--session-id", "status-running-test"}
    return nil
}
```
