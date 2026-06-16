# Scenario

**Feature**: a session directory exists with `meta.json` and `events.jsonl`

```
# implement agent reads requirement, writes code via Fake Codex
doctest agent implement --requirement req.md -> Fake Codex -> implementation

# session lifecycle
create session -> run sub-agent -> events recorded -> yield questions -> resume

# session id resolution order
--session-id flag -> opencode discovery -> codex resume -> fake-codex -> error
```

## Preconditions
- A session directory exists with `meta.json` and `events.jsonl`.
- The `--trace` flag follows an existing session and prints events.

## Steps
1. Create a session directory with explicit_session_id = `trace-test-sess`.
2. Write known events to `events.jsonl`.
3. Run `doctest agent implement --session-id trace-test-sess --trace`.

```go
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func Setup(t *testing.T, req *Request) error {
	sessHome := sessionsDir()
	now := time.Now()
	dateDir := now.Format("2006/01/02")
	sessName := fmt.Sprintf("sess_%s_%d", now.Format("150405"), now.UnixNano())
	sessDir := filepath.Join(sessHome, dateDir, sessName)
	os.MkdirAll(sessDir, 0755)

	meta := map[string]any{
		"explicit_session_id": "trace-test-sess",
		"created_at":          now.Format(time.RFC3339),
	}
	metaData, _ := json.MarshalIndent(meta, "", "  ")
	os.WriteFile(filepath.Join(sessDir, "meta.json"), metaData, 0644)

	events := []string{
		`{"type":"text","part":{"id":"m1","type":"message","text":"Hello from trace!"}}`,
	}
	eventsData := strings.Join(events, "\n") + "\n"
	os.WriteFile(filepath.Join(sessDir, "events.jsonl"), []byte(eventsData), 0644)

	req.Args = []string{"agent", "implement", "--session-id", "trace-test-sess", "--trace"}
	return nil
}
```
