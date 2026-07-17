---
label: heavy
---

## Expected
- First call: session directory created with `main_agent_codex_thread_id=flag-fallback-test` (no `explicit_session_id`).
- Second call: same session directory found and reused via `--session-id` codex fallback.
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

    sessDir, metaPath := getSessionDir(t, "main_agent_codex_thread_id", "flag-fallback-test")
    if sessDir == "" {
        t.Fatal("first call did not create a session dir with main_agent_codex_thread_id=flag-fallback-test")
    }

    m1 := readMetaJSON(t, metaPath)
    inner1, _ := m1["opencode_session_id"].(string)
    t.Logf("first call: session dir = %s, opencode_session_id = %s", sessDir, inner1)

    if _, ok := m1["explicit_session_id"]; ok {
        t.Fatalf("first call should not have explicit_session_id (created via CODEX_THREAD_ID), got %v", m1["explicit_session_id"])
    }

    if !strings.Contains(resp.Stdout, "first fallback call done") {
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
        "session_id":"inner-fallback-2",
        "llm_events":[
            {"type":"message","text":"second fallback call done"}
        ]
    }`), 0644)

    var fakeCodexPath string
    for _, e := range req.Env {
        if strings.HasPrefix(e, "AGENT_RUNNER_FAKE_CODEX_PATH=") {
            fakeCodexPath = e[len("AGENT_RUNNER_FAKE_CODEX_PATH="):]
            break
        }
    }

    cmd := exec.Command(doctestBin, "agent", "implement", "--session-id", "flag-fallback-test", "--agent-runner", "fake-codex", "second fallback call")
    cmd.Env = append(os.Environ(),
        "FAKE_CODEX_MOCK_CONFIG="+mockPath,
        "AGENT_RUNNER_FAKE_CODEX_PATH="+fakeCodexPath,
    )
    out2, err2 := cmd.CombinedOutput()
    if err2 != nil {
        t.Fatalf("second call failed: %v\noutput:\n%s", err2, string(out2))
    }
    t.Logf("second call output: %s", strings.TrimSpace(string(out2)))

    if !strings.Contains(string(out2), "second fallback call done") {
        t.Fatalf("second call stdout missing expected text:\n%s", string(out2))
    }

    sessDirExplicit, _ := getSessionDir(t, "explicit_session_id", "flag-fallback-test")
    if sessDirExplicit != "" && sessDirExplicit != sessDir {
        t.Fatalf("second call created a NEW session via explicit_session_id instead of reusing existing one.\n  existing (codex): %s\n  new (explicit): %s", sessDir, sessDirExplicit)
    }

    _, metaPath2 := getSessionDir(t, "main_agent_codex_thread_id", "flag-fallback-test")
    if metaPath2 == "" {
        t.Fatal("original session not found after second call")
    }
    if metaPath2 != metaPath {
        t.Fatalf("session meta path changed: %s → %s (session was not reused)", metaPath, metaPath2)
    }

    t.Logf("verified session reused via --session-id codex fallback")
}
```
