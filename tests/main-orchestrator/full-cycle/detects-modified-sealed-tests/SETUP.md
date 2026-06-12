## Preconditions
- A test tree has been written and sealed with `git add`.
- The sub-agent (or its hook) modifies a sealed test file.

## Steps
1. Create and seal a test tree.
2. Modify a sealed test file (simulating sub-agent change).
3. Run `git diff` to detect the modification.

```go
import (
    "os"
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

    // Simulate sub-agent modifying a sealed test
    assertPath := filepath.Join(treeDir, "basic", "ASSERT.md")
    orig, err := os.ReadFile(assertPath)
    if err != nil {
        t.Fatal(err)
    }
    modified := strings.Replace(string(orig), "Hello, world!", "Hi, world!", 1)
    os.WriteFile(assertPath, []byte(modified), 0644)

    req.Env = append(req.Env, "REPO_DIR="+repoDir)
    req.Args = []string{"--help"}
    return nil
}
```
