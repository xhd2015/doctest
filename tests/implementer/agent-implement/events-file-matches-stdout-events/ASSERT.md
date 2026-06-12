## Expected
- Exit code 0.
- `events.jsonl` contains events matching the mock config `stdout_events`.
- Event count matches.
- Event IDs match expected values.

```go
import (
    "encoding/json"
    "path/filepath"
    "testing"
    "time"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    sessionsDir := sessionsDir()

    today := time.Now().Format("2006/01/02")
    dateDir := filepath.Join(sessionsDir, today)

    sessDir := ""
    metaPath := findSessionMeta(t, "main_agent_codex_thread_id", "impl_test_events_match")
    if metaPath == "" {
        t.Fatal("no session found for impl_test_events_match")
    }
    sessDir = filepath.Dir(metaPath)
    if sessDir == "" || sessDir == dateDir {
        t.Fatalf("unexpected session dir: %s", sessDir)
    }

    events := readEventsJSONL(t, sessDir)
    if len(events) == 0 {
        t.Fatal("events.jsonl has no events")
    }

    expectedIDs := []string{"evt_a", "evt_a", "evt_b"}
    if len(events) != len(expectedIDs) {
        t.Fatalf("expected %d events, got %d", len(expectedIDs), len(events))
    }

    for i, ev := range events {
        item, ok := ev["item"].(map[string]any)
        if !ok {
            t.Fatalf("event %d has no item object", i+1)
        }
        id, _ := item["id"].(string)
        if id != expectedIDs[i] {
            t.Fatalf("event %d: expected id %q, got %q", i+1, expectedIDs[i], id)
        }
    }

    itemsJSON, _ := json.Marshal(events)
    t.Logf("events.jsonl content:\n%s", string(itemsJSON))
}
```
