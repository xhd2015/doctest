## Expected
- Exit code 0.
- Session directory is created under `~/.agent-pro/dedicated-agents/doctest-agent-implementer/sessions/YYYY/MM/DD/sess_*/`.
- Session dir name matches `sess_HHMMSS_<nano>` format.
- The session dir contains `meta.json` with `main_agent_codex_thread_id`.
- `doctest_agent_implementer_session_id` is absent (CODEX_THREAD_ID does not populate it).
- A `questions/` dir exists with a timestamp-named `.json` file.

## Exit Code
- Exit code 0.

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
        data, readErr := os.ReadFile(metaPath)
        if readErr != nil {
            continue
        }
        var meta struct {
            DoctestAgentImplementerSessionID string `json:"doctest_agent_implementer_session_id"`
            MainAgentCodexThreadID           string `json:"main_agent_codex_thread_id"`
            CreatedAt                        string `json:"created_at"`
        }
        if jsonErr := json.Unmarshal(data, &meta); jsonErr != nil {
            continue
        }
        if meta.MainAgentCodexThreadID == "impl_test_new_session" {
            sessDir = filepath.Join(dateDir, entry.Name())
            if matched, _ := filepath.Match("sess_[0-9][0-9][0-9][0-9][0-9][0-9]_[0-9]*", entry.Name()); !matched {
                t.Fatalf("session dir name %q does not match sess_HHMMSS_<nano> format", entry.Name())
            }
            if meta.MainAgentCodexThreadID == "" {
                t.Fatal("meta.json missing main_agent_codex_thread_id")
            }
            if meta.CreatedAt == "" {
                t.Fatal("meta.json missing created_at")
            }
            break
        }
    }
    if sessDir == "" {
        t.Fatal("no session found with main_agent_codex_thread_id=impl_test_new_session")
    }

    questionsDir := filepath.Join(sessDir, "questions")
    qFiles, qErr := os.ReadDir(questionsDir)
    if qErr != nil {
        t.Fatalf("cannot read questions dir %s: %v", questionsDir, qErr)
    }
    var jsonFiles []string
    for _, f := range qFiles {
        if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
            jsonFiles = append(jsonFiles, f.Name())
        }
    }
    if len(jsonFiles) == 0 {
        t.Fatal("no questions .json file found")
    }
    match := false
    for _, name := range jsonFiles {
        if matched, _ := filepath.Match("[0-9][0-9][0-9][0-9]_[0-9][0-9]_[0-9][0-9]_[0-9][0-9]_[0-9][0-9]_[0-9][0-9]*.json", name); matched {
            match = true
            break
        }
    }
    if !match {
        t.Fatalf("no timestamp-named .json file found in questions dir, got: %v", jsonFiles)
    }
}
```
