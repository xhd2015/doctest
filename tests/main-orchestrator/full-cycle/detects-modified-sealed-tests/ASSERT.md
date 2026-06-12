## Expected
- `git diff` on the test directory shows the modification.
- The unstaged diff contains the changed assertion.

```go
import (
    "strings"
    "testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    repoDir := ""
    for _, env := range req.Env {
        if strings.HasPrefix(env, "REPO_DIR=") {
            repoDir = env[len("REPO_DIR="):]
            break
        }
    }
    if repoDir == "" {
        t.Fatal("REPO_DIR not set by SETUP")
    }

    diffOut, _, _ := runCmd(t, repoDir, nil, "git", "diff", "tests/greet")
    if !strings.Contains(diffOut, "Hi, world!") {
        t.Fatalf("git diff should contain the modification:\n%s", diffOut)
    }
    if !strings.Contains(diffOut, "Hello, world!") {
        t.Fatalf("git diff should show original text as removed:\n%s", diffOut)
    }
}
```
