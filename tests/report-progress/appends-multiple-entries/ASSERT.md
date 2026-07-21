---
label: heavy
---

## Expected
- The progress file contains two lines after two invocations.
- Each line is valid JSON with `type` = `"progress"`.
- The two descriptions differ.

```go
import (
    "bytes"
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("first run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("first run exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    tf := os.Getenv("TEST_PROGRESS_FILE")
    if tf == "" {
        t.Fatal("TEST_PROGRESS_FILE not set")
    }

    // Second invocation with a different description
    rpBin := filepath.Join(t.TempDir(), "report-progress")
    if out, err := exec.Command("cp", req.Bin, rpBin).CombinedOutput(); err != nil {
        t.Fatalf("copy report-progress for second run: %v\n%s", err, string(out))
    }

    cmd := exec.Command(rpBin, "second step")
    cmd.Env = os.Environ()
    cmd.Env = append(cmd.Env, "PROGRESS_FILE="+tf)
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    runErr := cmd.Run()
    exitCode := 0
    if runErr != nil {
        if exitErr, ok := runErr.(*exec.ExitError); ok {
            exitCode = exitErr.ExitCode()
        }
    }
    if exitCode != 0 {
        t.Fatalf("second run exit code = %d, want 0\nstderr:\n%s", exitCode, stderr.String())
    }

    data, readErr := os.ReadFile(tf)
    if readErr != nil {
        t.Fatalf("read progress file: %v", readErr)
    }
    lines := strings.Split(strings.TrimSpace(string(data)), "\n")
    if len(lines) != 2 {
        t.Fatalf("expected 2 lines in progress file, got %d\ncontent:\n%s", len(lines), string(data))
    }

    var descs []string
    for i, line := range lines {
        var entry map[string]any
        if jsonErr := json.Unmarshal([]byte(line), &entry); jsonErr != nil {
            t.Fatalf("line %d: invalid JSON: %v\n%s", i+1, jsonErr, line)
        }
        if entry["type"] != "progress" {
            t.Fatalf("line %d: expected type=progress, got %v", i+1, entry["type"])
        }
        desc, _ := entry["description"].(string)
        if desc == "" {
            t.Fatalf("line %d: missing description", i+1)
        }
        descs = append(descs, desc)
        if entry["timestamp"] == nil {
            t.Fatalf("line %d: missing timestamp", i+1)
        }
    }

    if descs[0] != "first step" {
        t.Fatalf("first entry description = %q, want 'first step'", descs[0])
    }
    if descs[1] != "second step" {
        t.Fatalf("second entry description = %q, want 'second step'", descs[1])
    }
    if descs[0] == descs[1] {
        t.Fatal("expected different descriptions for two entries")
    }
    fmt.Fprint(os.Stderr, "two entries appended and valid JSONL")
}
```
