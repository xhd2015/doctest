## Preconditions
- `tests/implementer/SETUP.md` has built doctest and fake-codex binaries.
- Session directories are under `DOCTEST_DEBUG_SESSION_HOME`.

## Steps
1. Provide helpers for writing requirement files and reading messages.jsonl.

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	req.Env = append(req.Env, "TEST_GROUP=requirement")
	return nil
}

func writeRequirementFile(t *testing.T, req *Request, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "REQUIREMENT.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write requirement file: %v", err)
	}
	return path
}

func readMessagesFile(t *testing.T, sessionDir string) string {
	t.Helper()
	msgsPath := filepath.Join(sessionDir, "messages.jsonl")
	data, err := os.ReadFile(msgsPath)
	if err != nil {
		t.Fatalf("read messages.jsonl in %s: %v", sessionDir, err)
	}
	return string(data)
}

func sessionHasMessagesFile(t *testing.T, sessionDir string) bool {
	t.Helper()
	msgsPath := filepath.Join(sessionDir, "messages.jsonl")
	_, err := os.Stat(msgsPath)
	return err == nil
}
```
