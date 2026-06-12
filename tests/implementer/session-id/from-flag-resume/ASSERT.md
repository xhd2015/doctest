## Expected
- First call: session directory is created with `explicit_session_id=resume-flag-777`.
- Second call: same session directory is found and reused.
- Only ONE session dir exists for this ID.

```go
import (
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("first call exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    sessDir, metaPath := getSessionDir(t, "explicit_session_id", "resume-flag-777")
    if sessDir == "" {
        t.Fatal("first call did not create a session dir with explicit_session_id=resume-flag-777")
    }

    m1 := readMetaJSON(t, metaPath)
    inner1, _ := m1["opencode_session_id"].(string)
    t.Logf("first call: session dir = %s, opencode_session_id = %s", sessDir, inner1)

    if !strings.Contains(resp.Stdout, "first flag call done") {
        t.Fatalf("first call stdout missing expected text:\n%s", resp.Stdout)
    }

    doctestBin := ""
    for _, e := range req.Env {
        if strings.HasPrefix(e, "DOCTEST_BIN_FOR_RESUME=") {
            doctestBin = e[len("DOCTEST_BIN_FOR_RESUME="):]
            break
        }
    }
    if doctestBin == "" {
        t.Fatal("DOCTEST_BIN_FOR_RESUME not set")
    }

    mockPath := filepath.Join(t.TempDir(), "mock2.json")
    os.WriteFile(mockPath, []byte(`{
        "version":"agent-pro.fake-runner.v1",
        "runner":"fake-codex",
        "session_id":"inner-flag-resume-2",
        "stdout_events":[
            {"type":"item.completed","item":{"id":"m1","type":"message","text":"second flag call done","status":"completed"}}
        ]
    }`), 0644)

    var fakeCodexPath string
    for _, e := range req.Env {
        if strings.HasPrefix(e, "AGENT_RUNNER_FAKE_CODEX_PATH=") {
            fakeCodexPath = e[len("AGENT_RUNNER_FAKE_CODEX_PATH="):]
            break
        }
    }

    cmd := exec.Command(doctestBin, "agent", "implement", "--session-id", "resume-flag-777", "--agent-runner", "fake-codex", "second call")
    cmd.Env = append(os.Environ(),
        "FAKE_CODEX_MOCK_CONFIG="+mockPath,
        "AGENT_RUNNER_FAKE_CODEX_PATH="+fakeCodexPath,
    )
    out2, err2 := cmd.CombinedOutput()
    if err2 != nil {
        t.Fatalf("second call failed: %v\noutput:\n%s", err2, string(out2))
    }
    t.Logf("second call output: %s", strings.TrimSpace(string(out2)))

    if !strings.Contains(string(out2), "second flag call done") {
        t.Fatalf("second call stdout missing expected text:\n%s", string(out2))
    }

    count := countSessionDirs(t, "explicit_session_id", "resume-flag-777")
    if count != 1 {
        t.Fatalf("expected 1 session dir, got %d (session was not reused)", count)
    }

    // Verify opencode_session_id in meta.json is preserved (not overwritten on resume)
    _, metaPath2 := getSessionDir(t, "explicit_session_id", "resume-flag-777")
    if metaPath2 == "" {
        t.Fatal("session dir not found after second call")
    }
    m2 := readMetaJSON(t, metaPath2)
    inner2, _ := m2["opencode_session_id"].(string)

    if inner1 != "" && inner2 != inner1 {
        t.Fatalf("opencode_session_id changed across calls: %q → %q (should be preserved on resume)", inner1, inner2)
    }
}
```
