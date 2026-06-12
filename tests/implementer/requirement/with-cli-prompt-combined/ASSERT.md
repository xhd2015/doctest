## Expected
- Exit code 0.
- messages.jsonl contains both the requirement file content and the CLI prompt.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}

	sessionDir, _ := getSessionDir(t, "main_agent_codex_thread_id", "codex-tid-335")
	if sessionDir == "" {
		t.Fatal("no session found with main_agent_codex_thread_id=codex-tid-335")
	}

	content := readMessagesFile(t, sessionDir)
	if !strings.Contains(content, "spec content") {
		t.Fatalf("messages.jsonl should contain 'spec content', got: %s", content)
	}
	if !strings.Contains(content, "question") {
		t.Fatalf("messages.jsonl should contain 'question', got: %s", content)
	}
}
```
