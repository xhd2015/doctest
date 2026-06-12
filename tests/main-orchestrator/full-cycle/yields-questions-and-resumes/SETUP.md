## Preconditions
- A stub test tree is written, RED confirmed, and sealed.
- The mock config has a `before_exit` hook that calls yield-pending-questions.
- After questions, the same session is resumed with answers.

## Steps
1. Create and seal stub test tree.
2. Run `doctest agent implement` with mock that yields questions.
3. Capture CODEX_THREAD_ID from output or env.
4. Provide answers via a second `doctest agent implement` call.

```go
import (
    "fmt"
    "os"
    "path/filepath"
    "testing"
)

func Setup(t *testing.T, req *Request) error {
    repoDir := filepath.Join(t.TempDir(), "repo")
    os.MkdirAll(repoDir, 0755)
    runCmd(t, repoDir, nil, "git", "init")
    runCmd(t, repoDir, nil, "git", "config", "user.email", "test@test.com")
    runCmd(t, repoDir, nil, "git", "config", "user.name", "Test")

    treeDir := filepath.Join(repoDir, "tests", "greet")
    createDoctestTree(t, treeDir, true)
    runCmd(t, repoDir, nil, "git", "add", "tests/greet")

    yieldPQ := os.Getenv("YIELD_PQ_BIN")
    if yieldPQ == "" {
        t.Fatal("YIELD_PQ_BIN not set")
    }

    hookCmd := yieldPQ + fmt.Sprintf(` '{"id":"1","question":"What should Greet(\"\") return?"}'`)
    writeMockConfig(t, req, fmt.Sprintf(`{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","session_id":"sess_yield","hook_command":%q,"hooks":[{"at":"before_exit","event":"yield","payload":{"ok":true}}],"stdout_events":[{"type":"item.completed","item":{"id":"m1","type":"message","text":"need clarification","status":"completed"}}]}`, hookCmd+" {{event}}"))

    req.Env = append(req.Env, "REPO_DIR="+repoDir)
    req.Env = append(req.Env, "CODEX_THREAD_ID=impl_test_orch_yield")
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "implement greet feature"}
    return nil
}
```
