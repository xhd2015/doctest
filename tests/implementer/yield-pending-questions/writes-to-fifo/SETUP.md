## Preconditions
- `QUESTION_FIFO` env var points to a writable file.
- The binary receives a question JSON argument.

## Steps
1. Set `QUESTION_FIFO` to a temp file.
2. Invoke yield-pending-questions with JSON.
3. Read the file and verify the question was written.

```go
import (
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    qFile := filepath.Join(t.TempDir(), "questions.jsonl")
    req.Env = append(req.Env, "QUESTION_FIFO="+qFile)
    os.Setenv("TEST_Q_FILE", qFile)
    req.Args = []string{`{"id":"1","question":"What is the target port?"}`}
    return nil
}
```
