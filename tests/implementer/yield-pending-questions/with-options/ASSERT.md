## Expected
- Exit code 0.
- The file contains a JSON line with `type` = `"question"` and includes options with explanations.

## Exit Code
- Exit code 0.

```go
import (
    "encoding/json"
    "os"
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatalf("run failed: %v", err)
    }
    if resp.ExitCode != 0 {
        t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
    }

    qFile := os.Getenv("TEST_Q_FILE")
    if qFile == "" {
        t.Fatal("TEST_Q_FILE not set")
    }

    data, readErr := os.ReadFile(qFile)
    if readErr != nil {
        t.Fatalf("read questions file: %v", readErr)
    }
    content := string(data)
    if content == "" {
        t.Fatal("questions file is empty")
    }

    var entry struct {
        Type     string `json:"type"`
        ID       string `json:"id"`
        Question string `json:"question"`
        Options  []struct {
            Option      string `json:"option"`
            Explanation string `json:"explanation"`
        } `json:"options"`
    }
    line := strings.TrimSpace(content)
    if jsonErr := json.Unmarshal([]byte(line), &entry); jsonErr != nil {
        t.Fatalf("invalid JSON: %v\n%s", jsonErr, line)
    }
    if entry.Type != "question" {
        t.Fatalf("expected type 'question', got %q", entry.Type)
    }
    if entry.Question != "What is the target port?" {
        t.Fatalf("expected question text, got %q", entry.Question)
    }
    if len(entry.Options) != 2 {
        t.Fatalf("expected 2 options, got %d", len(entry.Options))
    }
    if entry.Options[0].Option != "3000" {
        t.Fatalf("expected first option '3000', got %q", entry.Options[0].Option)
    }
    if entry.Options[0].Explanation == "" {
        t.Fatal("expected first option to have explanation")
    }
    if entry.Options[1].Option != "8080" {
        t.Fatalf("expected second option '8080', got %q", entry.Options[1].Option)
    }
    if entry.Options[1].Explanation == "" {
        t.Fatal("expected second option to have explanation")
    }
}
```
