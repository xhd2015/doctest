## Expected
- Exit code 0.
- messages.jsonl contains a "create_time" field.
- "create_time" is a valid RFC 3339 datetime string.

```go
import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}

	sessionDir, _ := getSessionDir(t, "main_agent_codex_thread_id", "codex-tid-338")
	if sessionDir == "" {
		t.Fatal("no session found with main_agent_codex_thread_id=codex-tid-338")
	}

	content := readMessagesFile(t, sessionDir)
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		t.Fatal("messages.jsonl is empty")
	}

	var entry map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &entry); err != nil {
		t.Fatalf("failed to parse messages.jsonl entry: %v\ncontent: %s", err, content)
	}

	createTime, ok := entry["create_time"].(string)
	if !ok || createTime == "" {
		t.Fatalf("create_time missing or empty, entry: %v", entry)
	}

	if _, err := time.Parse(time.RFC3339, createTime); err != nil {
		t.Fatalf("create_time is not valid RFC 3339: %s, error: %v", createTime, err)
	}
}
```
