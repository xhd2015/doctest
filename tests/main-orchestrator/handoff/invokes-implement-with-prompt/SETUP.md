## Preconditions
- A test tree exists with stub Run().
- Tests are confirmed RED and sealed.
- A mock config is ready for the sub-agent.

## Steps
1. Create a doctest tree and seal it.
2. Write a mock config with a completion event.
3. Run `doctest agent implement "implement greet feature" --agent-runner fake-codex`.

```go
import (
    "os"
    "os/exec"
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
    createDoctestTree(t, treeDir, false)
    runCmd(t, repoDir, nil, "git", "add", "tests/greet")

    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","session_id":"sess_handoff","stdout_events":[{"type":"item.completed","item":{"id":"m1","type":"message","text":"implemented greet feature","status":"completed"}}]}`)

    req.Env = append(req.Env, "REPO_DIR="+repoDir)
    req.Env = append(req.Env, "CODEX_THREAD_ID=impl_test_orch_handoff")
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "implement greet feature"}
    _ = exec.Command
    return nil
}
```
