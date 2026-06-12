## Expected
- `git diff --cached --name-only` lists all test files.
- `git diff --name-only` (unstaged) is empty.

```go
import (
    "path/filepath"
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

    stagedOut, stagedErr, _ := runCmd(t, repoDir, nil, "git", "diff", "--cached", "--name-only")
    _ = stagedErr
    stagedFiles := strings.TrimSpace(stagedOut)
    if !strings.Contains(stagedFiles, "DOCTEST.md") {
        t.Fatalf("DOCTEST.md not staged:\n%s", stagedOut)
    }
    if !strings.Contains(stagedFiles, filepath.Join("greet", "SETUP.md")) {
        t.Fatalf("greet/SETUP.md not staged:\n%s", stagedOut)
    }

    unsOut, _, _ := runCmd(t, repoDir, nil, "git", "diff", "--name-only")
    unsFiles := strings.TrimSpace(unsOut)
    if unsFiles != "" {
        t.Fatalf("unexpected unstaged changes:\n%s", unsOut)
    }
}
```
