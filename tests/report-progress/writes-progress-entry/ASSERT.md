---
label: e2e, heavy
---

## Expected
- Exit code 0.
- The file contains a single JSON line with `type` = `"progress"`.
- The line contains the description text and a valid timestamp.

```go
import (
    "encoding/json"
    "os"
    "strings"
    "testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    tf := req.ProgressFile
	if tf == "" {
		t.Fatal("req.ProgressFile not set")
	}

    data, readErr := os.ReadFile(tf)
    if readErr != nil {
        t.Fatalf("read progress file: %v", readErr)
    }
    content := strings.TrimSpace(string(data))
    if content == "" {
        t.Fatal("progress file is empty")
    }

    var entry map[string]any
    if jsonErr := json.Unmarshal([]byte(content), &entry); jsonErr != nil {
        t.Fatalf("parse progress entry: %v\ncontent:\n%s", jsonErr, content)
    }

    if entry["type"] != "progress" {
        t.Fatalf("expected type=progress, got %v\ncontent:\n%s", entry["type"], content)
    }
    desc, _ := entry["description"].(string)
    if desc != "implementing JSON parser" {
        t.Fatalf("expected description 'implementing JSON parser', got %q", desc)
    }
    ts, _ := entry["timestamp"].(string)
    if ts == "" {
        t.Fatal("timestamp is missing or empty")
    }
}
```
