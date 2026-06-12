## Expected
- `doctest_agent_implementer_session_id`: `store-prio-1`
- `explicit_session_id`: absent (flag not used)
- `main_agent_codex_thread_id`: `store-codex-2`
- `main_agent_opencode_session_id`: absent (env var took priority, discovery not used)
- `opencode_session_id`: non-empty (inner spawned agent session)
- `created_at`: non-empty

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    metaPath := findSessionMeta(t, "doctest_agent_implementer_session_id", "store-prio-1")
    if metaPath == "" {
        t.Fatal("no session found with doctest_agent_implementer_session_id=store-prio-1")
    }

    m := readMetaJSON(t, metaPath)

    gotFields := make(map[string]bool)
    for k := range m {
        gotFields[k] = true
    }

    checkField := func(field string, expected string) {
        v, ok := m[field].(string)
        if !ok || v != expected {
            t.Fatalf("%s = %q, want %q", field, v, expected)
        }
    }

    checkFieldPresent := func(field string) {
        v, ok := m[field].(string)
        if !ok || v == "" {
            t.Fatalf("%s must be set and non-empty, got %v", field, m[field])
        }
    }

    checkField("doctest_agent_implementer_session_id", "store-prio-1")
    checkField("main_agent_codex_thread_id", "store-codex-2")
    checkFieldPresent("opencode_session_id")
    checkFieldPresent("created_at")

    if _, ok := m["explicit_session_id"]; ok {
        t.Fatalf("explicit_session_id should NOT be present when --session-id flag not used, got %v", m["explicit_session_id"])
    }
    if _, ok := m["main_agent_opencode_session_id"]; ok {
        t.Fatalf("main_agent_opencode_session_id should NOT be present when env var takes priority, got %v", m["main_agent_opencode_session_id"])
    }

    if gotFields["codex_thread_id"] {
        t.Fatal("old field codex_thread_id must not appear in meta.json")
    }
}
```
