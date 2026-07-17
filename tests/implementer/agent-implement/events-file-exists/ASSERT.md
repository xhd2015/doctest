---
label: heavy
---

## Expected
- Exit code 0.
- Session directory is created.
- `events.jsonl` exists in the session directory and is a regular file with content.

```go
import (
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

    var sessDir string
    for _, entry := range entries {
        if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "sess_") {
            continue
        }
        metaPath := filepath.Join(dateDir, entry.Name(), "meta.json")
        m := readMetaJSON(t, metaPath)
        if v, ok := m["main_agent_codex_thread_id"]; ok {
            if s, ok := v.(string); ok && s == "impl_test_events_exists" {
                sessDir = filepath.Join(dateDir, entry.Name())
                break
            }
        }
    }
    if sessDir == "" {
        t.Fatal("no session found for impl_test_events_exists")
    }

    eventsPath := filepath.Join(sessDir, "events.jsonl")
    info, statErr := os.Stat(eventsPath)
    if statErr != nil {
        t.Fatalf("events.jsonl does not exist: %v", statErr)
    }
    if info.IsDir() {
        t.Fatal("events.jsonl is a directory, expected a regular file")
    }
    if info.Size() == 0 {
        t.Fatal("events.jsonl is empty")
    }
}
```
