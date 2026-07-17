---
label: heavy
---

## Expected
- Exit code 0.
- `events.jsonl` contains both the pre-existing initial event and the new resume event.
- The initial event (id `initial_evt`) appears before the resume event (id `resume_evt`).

```go
import (
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
    sessDir := filepath.Join(dateDir, "sess_test_events_append")

    events := readEventsJSONL(t, sessDir)
    if len(events) < 2 {
        t.Fatalf("expected at least 2 events (initial + resume), got %d", len(events))
    }

    foundInitial := false
    foundResume := false
    initialIdx := -1
    resumeIdx := -1
    for i, ev := range events {
        item, ok := ev["item"].(map[string]any)
        if !ok {
            continue
        }
        id, _ := item["id"].(string)
        if id == "initial_evt" {
            foundInitial = true
            initialIdx = i
        }
        if id == "resume_evt" {
            foundResume = true
            resumeIdx = i
        }
    }

    if !foundInitial {
        t.Fatal("initial event (id=initial_evt) not found in events.jsonl")
    }
    if !foundResume {
        t.Fatal("resume event (id=resume_evt) not found in events.jsonl")
    }
    if initialIdx >= resumeIdx {
        t.Fatalf("expected initial event before resume event, got initial at %d and resume at %d", initialIdx, resumeIdx)
    }
}
```
