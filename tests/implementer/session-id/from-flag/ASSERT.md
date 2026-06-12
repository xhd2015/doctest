## Expected
- Exit code 0.
- `explicit_session_id` = `my-sess-flag` (from `--session-id` flag).
- `doctest_agent_implementer_session_id` is ABSENT (env var not set).
- `main_agent_codex_thread_id` = `should-be-ignored` (stored for traceability).
- `opencode_session_id` is set (inner agent session).
- `created_at` is set.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    metaPath := findSessionMeta(t, "explicit_session_id", "my-sess-flag")
    if metaPath == "" {
        t.Fatal("no session found with explicit_session_id=my-sess-flag")
    }

    m := readMetaJSON(t, metaPath)

    if v, ok := m["explicit_session_id"].(string); !ok || v != "my-sess-flag" {
        t.Fatalf("explicit_session_id = %v, want my-sess-flag", m["explicit_session_id"])
    }
    if _, ok := m["doctest_agent_implementer_session_id"]; ok {
        t.Fatalf("doctest_agent_implementer_session_id should NOT be present when using --session-id flag, got %v", m["doctest_agent_implementer_session_id"])
    }
    if v, ok := m["main_agent_codex_thread_id"].(string); !ok || v != "should-be-ignored" {
        t.Fatalf("main_agent_codex_thread_id = %v, want should-be-ignored", m["main_agent_codex_thread_id"])
    }
    if _, ok := m["main_agent_opencode_session_id"]; ok {
        t.Fatalf("main_agent_opencode_session_id should NOT be present when flag is set, got %v", m["main_agent_opencode_session_id"])
    }
    if v, ok := m["opencode_session_id"].(string); !ok || v == "" {
        t.Fatalf("opencode_session_id must be set, got %v", m["opencode_session_id"])
    }
    if _, ok := m["created_at"]; !ok {
        t.Fatal("created_at must be set")
    }

    t.Logf("stdout: %s", strings.TrimSpace(resp.Stdout))
}
```
