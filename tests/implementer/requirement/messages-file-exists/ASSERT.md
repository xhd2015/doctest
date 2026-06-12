## Expected
- Exit code 0.
- messages.jsonl exists in the session directory.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}

	sessionDir, _ := getSessionDir(t, "main_agent_codex_thread_id", "codex-tid-333")
	if sessionDir == "" {
		t.Fatal("no session found with main_agent_codex_thread_id=codex-tid-333")
	}

	if !sessionHasMessagesFile(t, sessionDir) {
		t.Fatalf("messages.jsonl not found in session dir: %s", sessionDir)
	}
}
```
