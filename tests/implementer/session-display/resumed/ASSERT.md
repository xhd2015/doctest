## Expected
- First call succeeds and shows "Session created: sess-display-resume-1".
- Second call with same session ID shows "Session resumed: sess-display-resume-1".

```go
import (
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("first call failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("first call exit code = %d, want 0\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
    }
    if !strings.Contains(resp.Stdout, "Session created: sess-display-resume-1") {
        t.Fatalf("first call: expected stdout to contain 'Session created: sess-display-resume-1', got:\n%s", resp.Stdout)
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
        "session_id":"inner-sess-resume-2",
        "llm_events":[
            {"type":"message","text":"second call done"}
        ]
    }`), 0644)

    var fakeCodexPath string
    for _, e := range req.Env {
        if strings.HasPrefix(e, "AGENT_RUNNER_FAKE_CODEX_PATH=") {
            fakeCodexPath = e[len("AGENT_RUNNER_FAKE_CODEX_PATH="):]
            break
        }
    }

    cmd := exec.Command(doctestBin, "agent", "implement", "--session-id", "sess-display-resume-1", "--agent-runner", "fake-codex", "second")
    cmd.Env = append(os.Environ(),
        "FAKE_CODEX_MOCK_CONFIG="+mockPath,
        "AGENT_RUNNER_FAKE_CODEX_PATH="+fakeCodexPath,
    )
    out2, err2 := cmd.CombinedOutput()
    if err2 != nil {
        t.Fatalf("second call failed: %v\noutput:\n%s", err2, string(out2))
    }

    if !strings.Contains(string(out2), "Session resumed: sess-display-resume-1") {
        t.Fatalf("second call: expected stdout to contain 'Session resumed: sess-display-resume-1', got:\n%s", string(out2))
    }
}
```
