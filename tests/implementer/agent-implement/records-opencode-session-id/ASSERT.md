## Expected
- Exit code 0.
- Session directory is created with `meta.json`.
- `meta.json` contains `main_agent_codex_thread_id` and `opencode_session_id`.
- `opencode_session_id` is non-empty.

```go
import (
    "encoding/json"
    "os"
    "path/filepath"
    "strings"
    "testing"
    "time"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    sessionsDir := sessionsDir()

    today := time.Now().Format("2006/01/02")
    dateDir := filepath.Join(sessionsDir, today)
    entries, readErr := os.ReadDir(dateDir)
    if readErr != nil {
        t.Fatalf("cannot read date dir %s: %v", dateDir, readErr)
    }

    var newestDir string
    var newestTime int64
    for _, entry := range entries {
        if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "sess_") {
            continue
        }
        info, statErr := entry.Info()
        if statErr != nil {
            continue
        }
        modTime := info.ModTime().UnixNano()
        if modTime > newestTime {
            newestTime = modTime
            newestDir = entry.Name()
        }
    }
    if newestDir == "" {
        t.Fatal("no session directory created")
    }

    sessDir := filepath.Join(dateDir, newestDir)

    metaPath := filepath.Join(sessDir, "meta.json")
    data, readErr := os.ReadFile(metaPath)
    if readErr != nil {
        t.Fatalf("cannot read meta.json: %v", readErr)
    }
    var meta struct {
        MainAgentCodexThreadID string `json:"main_agent_codex_thread_id"`
        OpencodeSessionID      string `json:"opencode_session_id"`
        CreatedAt              string `json:"created_at"`
    }
    if jsonErr := json.Unmarshal(data, &meta); jsonErr != nil {
        t.Fatalf("invalid meta.json: %v", jsonErr)
    }
    if meta.MainAgentCodexThreadID == "" {
        t.Fatal("meta.json missing main_agent_codex_thread_id")
    }
    if meta.OpencodeSessionID == "" {
        t.Fatal("meta.json missing opencode_session_id")
    }
    if meta.CreatedAt == "" {
        t.Fatal("meta.json missing created_at")
    }
}
```
