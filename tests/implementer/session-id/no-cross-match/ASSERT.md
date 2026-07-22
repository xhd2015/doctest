---
label: heavy
---

## Expected
- First call: session created with `explicit_session_id=sess-A`.
- Second call: NEW session created with `main_agent_codex_thread_id=sess-B`.
- TWO session dirs exist (no cross-matching between fields).

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

    if !strings.Contains(resp.Stdout, "first call A done") {
        t.Fatalf("first call stdout missing expected text:\n%s", resp.Stdout)
    }

    // Verify the first session exists with explicit_session_id
    sessDirA, _ := getSessionDir(t, req, "explicit_session_id", "sess-A")
    if sessDirA == "" {
        t.Fatal("first call did not create session with explicit_session_id=sess-A")
    }
    t.Logf("first session dir: %s", sessDirA)

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
        "session_id":"inner-cross-B",
        "llm_events":[
            {"type":"message","text":"second call B done"}
        ]
    }`), 0644)

    var fakeCodexPath string
    for _, e := range req.Env {
        if strings.HasPrefix(e, "AGENT_RUNNER_FAKE_CODEX_PATH=") {
            fakeCodexPath = e[len("AGENT_RUNNER_FAKE_CODEX_PATH="):]
            break
        }
    }

    // Second call must use the same isolated session home as the first call
    // (DOCTEST_DEBUG_SESSION_HOME from root Setup). Without it, the product
    // looks under the process default session base and either pollutes or
    // resumes a leftover session — never writing sess-B into req.SessionHome.
    sessionHome := req.SessionHome
    if sessionHome == "" {
        for _, e := range req.Env {
            if strings.HasPrefix(e, "DOCTEST_DEBUG_SESSION_HOME=") {
                sessionHome = e[len("DOCTEST_DEBUG_SESSION_HOME="):]
                break
            }
        }
    }
    if sessionHome == "" {
        t.Fatal("DOCTEST_DEBUG_SESSION_HOME / req.SessionHome not set")
    }

    // Second call: use CODEX_THREAD_ID instead of --session-id
    cmd := exec.Command(doctestBin, "agent", "implement", "--agent-runner", "fake-codex", "second call")
    cmd.Env = append(os.Environ(),
        "CODEX_THREAD_ID=sess-B",
        "FAKE_CODEX_MOCK_CONFIG="+mockPath,
        "AGENT_RUNNER_FAKE_CODEX_PATH="+fakeCodexPath,
        "DOCTEST_DEBUG_SESSION_HOME="+sessionHome,
    )
    out2, err2 := cmd.CombinedOutput()
    if err2 != nil {
        t.Fatalf("second call failed: %v\noutput:\n%s", err2, string(out2))
    }
    t.Logf("second call output: %s", strings.TrimSpace(string(out2)))

    if !strings.Contains(string(out2), "second call B done") {
        t.Fatalf("second call stdout missing expected text:\n%s", string(out2))
    }
    // New source ID must create, not resume sess-A or any other session.
    if !strings.Contains(string(out2), "Session created: sess-B") {
        t.Fatalf("second call expected 'Session created: sess-B' (no cross-match), got:\n%s", string(out2))
    }

    // Verify TWO separate sessions exist
    if findSessionMeta(t, req, "explicit_session_id", "sess-A") == "" {
        t.Fatal("first session with explicit_session_id=sess-A was lost")
    }
    if findSessionMeta(t, req, "main_agent_codex_thread_id", "sess-B") == "" {
        t.Fatal("second session with main_agent_codex_thread_id=sess-B not found (should not match explicit_session_id)")
    }

    // First session should still only have explicit_session_id=sess-A, not codex
    metaA := readMetaJSON(t, findSessionMeta(t, req, "explicit_session_id", "sess-A"))
    if _, ok := metaA["main_agent_codex_thread_id"]; ok {
        t.Fatalf("first session should not have main_agent_codex_thread_id, got %v", metaA["main_agent_codex_thread_id"])
    }

    // Second session should have main_agent_codex_thread_id=sess-B
    metaB := readMetaJSON(t, findSessionMeta(t, req, "main_agent_codex_thread_id", "sess-B"))
    if v, _ := metaB["main_agent_codex_thread_id"].(string); v != "sess-B" {
        t.Fatalf("second session main_agent_codex_thread_id = %v, want sess-B", metaB["main_agent_codex_thread_id"])
    }
    if _, ok := metaB["explicit_session_id"]; ok {
        t.Fatalf("second session should not have explicit_session_id, got %v", metaB["explicit_session_id"])
    }
}
```
