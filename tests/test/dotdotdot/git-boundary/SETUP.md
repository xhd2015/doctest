# Scenario

**Feature**: tests for git boundary

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Group: git-boundary
Tests for `./...` respecting git repository boundaries.

When walking up to find a module root or walking down to discover tests,
`doctest` must stop at git repository boundaries (`.git` directories).

```go
import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "testing"
)

func initGitRepo(dir string) error {
    if out, err := exec.Command("git", "-C", dir, "init").CombinedOutput(); err != nil {
        return fmt.Errorf("git init: %v\n%s", err, string(out))
    }
    return nil
}

func writeGoMod(dir, modulePath string) error {
    content := "module " + modulePath + "\ngo 1.21\n"
    return os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0644)
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    t.Logf("git-boundary group: WorkDir=%s", req.WorkDir)
    return nil
}
```
