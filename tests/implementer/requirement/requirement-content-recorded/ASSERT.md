## Expected
- Exit code 0.
- messages.jsonl contains the requirement file content "feature spec here".

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}

	sessionDir, _ := getSessionDir(t, "main_agent_codex_thread_id", "codex-tid-334")
	if sessionDir == "" {
		t.Fatal("no session found with main_agent_codex_thread_id=codex-tid-334")
	}

	content := readMessagesFile(t, sessionDir)
	if !strings.Contains(content, "feature spec here") {
		t.Fatalf("messages.jsonl should contain 'feature spec here', got: %s", content)
	}
}
```
