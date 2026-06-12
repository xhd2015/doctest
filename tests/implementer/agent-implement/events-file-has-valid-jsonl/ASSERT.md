## Expected
- Exit code 0.
- `events.jsonl` exists in the session directory.
- Every line in `events.jsonl` is valid JSON.

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

    var sessDir string
    for _, entry := range entries {
        if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "sess_") {
            continue
        }
        metaPath := filepath.Join(dateDir, entry.Name(), "meta.json")
        m := readMetaJSON(t, metaPath)
        if v, ok := m["main_agent_codex_thread_id"]; ok {
            if s, ok := v.(string); ok && s == "impl_test_events_valid" {
                sessDir = filepath.Join(dateDir, entry.Name())
                break
            }
        }
    }
    if sessDir == "" {
        t.Fatal("no session found for impl_test_events_valid")
    }

    eventsPath := filepath.Join(sessDir, "events.jsonl")
    data, readErr := os.ReadFile(eventsPath)
    if readErr != nil {
        t.Fatalf("cannot read events.jsonl: %v", readErr)
    }

    lines := strings.Split(strings.TrimSpace(string(data)), "\n")
    if len(lines) == 0 {
        t.Fatal("events.jsonl has no lines")
    }

    for i, line := range lines {
        if strings.TrimSpace(line) == "" {
            continue
        }
        var ev map[string]any
        if jsonErr := json.Unmarshal([]byte(line), &ev); jsonErr != nil {
            t.Fatalf("line %d is not valid JSON: %v\nline: %s", i+1, jsonErr, line)
        }
    }
}
```
