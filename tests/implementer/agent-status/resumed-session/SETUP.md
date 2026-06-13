## Preconditions
- A session directory exists with `meta.json` from a previous run but NO `pid` file (previous session finished).
- A resumed session is started in the background using a slow mock config.
- `--status` should show "running" while the resumed session is still alive.

## Steps
1. Create a session directory with `meta.json` (explicit_session_id), no `pid` file.
2. Write mock config with `delay_ms` to keep the background process alive.
3. Start `doctest agent implement --session-id <ID> --agent-runner fake-codex "resumed task"` in the background.
4. Wait for the PID file to appear (poll with timeout).
5. Run `doctest agent implement --status --session-id <ID>`.

```go
import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
    "time"
)

func Setup(t *testing.T, req *Request) error {
    sessID := "status-resumed-test"
    sessHome := sessionsDir()

    now := time.Now()
    dateDir := now.Format("2006/01/02")
    sessName := fmt.Sprintf("sess_%s_%d", now.Format("150405"), now.UnixNano())
    sessDir := filepath.Join(sessHome, dateDir, sessName)

    if err := os.MkdirAll(sessDir, 0755); err != nil {
        t.Fatalf("create session dir: %v", err)
    }

    meta := map[string]any{
        "explicit_session_id": sessID,
        "agent_runner":        "fake-codex",
        "created_at":          now.Format(time.RFC3339),
    }
    metaData, _ := json.MarshalIndent(meta, "", "  ")
    if err := os.WriteFile(filepath.Join(sessDir, "meta.json"), metaData, 0644); err != nil {
        t.Fatalf("write meta.json: %v", err)
    }

    mockConfig := fmt.Sprintf(`{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","delay_ms":10000,"session_id":"%s","llm_events":[{"type":"message","text":"resumed and done"}]}`, sessID)
    mockPath := filepath.Join(t.TempDir(), "mock-resume.json")
    if err := os.WriteFile(mockPath, []byte(mockConfig), 0644); err != nil {
        t.Fatalf("write mock config: %v", err)
    }

    bin := req.Bin
    if bin == "" {
        t.Fatal("req.Bin is empty")
    }

    bgCmd := exec.Command(bin, "agent", "implement", "--agent-runner", "fake-codex", "--session-id", sessID, "resumed task")
    bgCmd.Env = append(os.Environ(),
        "FAKE_CODEX_MOCK_CONFIG="+mockPath,
        "DOCTEST_DEBUG_SESSION_HOME="+sessHome,
    )

    var bgStdout, bgStderr strings.Builder
    bgCmd.Stdout = &bgStdout
    bgCmd.Stderr = &bgStderr

    if err := bgCmd.Start(); err != nil {
        t.Fatalf("start background agent: %v", err)
    }

    defers := &struct{}{}
    _ = defers
    t.Cleanup(func() {
        if bgCmd.Process != nil {
            bgCmd.Process.Kill()
            bgCmd.Wait()
        }
    })

    pidPath := filepath.Join(sessDir, "pid")
    found := false
    for i := 0; i < 50; i++ {
        if _, err := os.Stat(pidPath); err == nil {
            found = true
            break
        }
        time.Sleep(100 * time.Millisecond)
    }

    if !found {
        t.Logf("bg stdout:\n%s", bgStdout.String())
        t.Logf("bg stderr:\n%s", bgStderr.String())
        t.Fatal("PID file did not appear within 5s — writeSessionPID not called for resumed session")
    }

    req.Args = []string{"agent", "implement", "--status", "--session-id", sessID}
    return nil
}
```
