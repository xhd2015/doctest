## Preconditions
- The repo supports git.
- A test tree exists in a temp dir.

## Steps
1. git init in a temp dir.
2. Create a doctest tree.
3. git add the test directory.
4. Verify files are staged.

```go
import (
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    repoDir := filepath.Join(t.TempDir(), "repo")
    os.MkdirAll(repoDir, 0755)
    runCmd(t, repoDir, nil, "git", "init")
    runCmd(t, repoDir, nil, "git", "config", "user.email", "test@test.com")
    runCmd(t, repoDir, nil, "git", "config", "user.name", "Test")

    treeDir := filepath.Join(repoDir, "tests", "greet")
    createDoctestTree(t, treeDir, false)

    runCmd(t, repoDir, nil, "git", "add", "tests/greet")

    req.Env = append(req.Env, "REPO_DIR="+repoDir)
    req.Args = []string{"--help"}
    _ = strings.Contains
    return nil
}
```
