## Preconditions
- A stub doctest tree exists.
- `doctest test` confirms RED.
- `git add` seals the tests.
- `doctest agent implement` hands off to sub-agent.

## Steps
1. Create stub test tree.
2. Run `doctest test` — verify RED.
3. `git add` — seal tests.
4. Run `doctest agent implement` with mock config — sub-agent completes.
5. Verify all phases succeeded.

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
    runCmd(t, repoDir, nil, "go", "mod", "init", "test")
    runCmd(t, repoDir, nil, "git", "init")
    runCmd(t, repoDir, nil, "git", "config", "user.email", "test@test.com")
    runCmd(t, repoDir, nil, "git", "config", "user.name", "Test")

    treeDir := filepath.Join(repoDir, "tests", "greet")
    createDoctestTree(t, treeDir, true) // stub

    // Phase 3: RED — run doctest from repo root (has go.mod), use absolute path to tree
    doctestBin := ""
    for _, env := range req.Env {
        if len(env) > 12 && env[:12] == "DOCTEST_BIN=" {
            doctestBin = env[12:]
            break
        }
    }
    out, errOut, code := runCmd(t, req.WorkDir, req.Env, doctestBin, "test", "-v", treeDir)
    allOut := out + errOut
    if code == 0 {
        t.Fatalf("expected RED, got exit 0\n%s", allOut)
    }
    if !strings.Contains(allOut, "not implemented") {
        t.Fatalf("expected 'not implemented' in RED output:\nSTDOUT:\n%s\nSTDERR:\n%s", out, errOut)
    }

    // Phase 4: Seal
    runCmd(t, repoDir, nil, "git", "add", "tests/greet")
    stagedOut, _, _ := runCmd(t, repoDir, nil, "git", "diff", "--cached", "--name-only")
    if !strings.Contains(stagedOut, "DOCTEST.md") {
        t.Fatalf("tests not staged after git add:\n%s", stagedOut)
    }

    // Phase 5: Handoff
    writeMockConfig(t, req, `{"version":"agent-pro.fake-runner.v1","runner":"fake-codex","session_id":"sess_full","stdout_events":[{"type":"item.completed","item":{"id":"m1","type":"message","text":"feature implemented","status":"completed"}}]}`)
    req.Env = append(req.Env, "REPO_DIR="+repoDir)
    req.Env = append(req.Env, "CODEX_THREAD_ID=impl_test_orch_full")
    req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex", "implement greet feature"}
    return nil
}
```
