# Cross-Module Git-Aware Discovery Tests

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **system under test** — behavior covered by this tree.

### Behaviors
- **run** — executes the scenarios in this suite.


Integration tests for `doctest test ./...` when a nested `go.mod` boundary is
crossed. At non-child module paths, discovery depends on whether ancestor and
nested module roots share the same git work tree.

## Git Context Notes

`gitRoot(dir)` resolves via `git.ShowToplevel(dir)` when inside a work tree,
otherwise `""`.

**Collapsed case:** A child directory nested inside the parent's git worktree
without its own `.git` is **same-git context** — `gitRoot(child) ==
gitRoot(parent)`. This is identical to `parent-git-child-same-git-discovers`;
there is no separate warn path for "parent in git, child has no own repo" when
the child lives inside the parent worktree.

**Distinct mismatch cases** require unequal git roots: separate repos
(`parent-git-child-other-git-warns`) or one side in git and the other outside
any work tree (`parent-not-git-child-git-warns`).

## Decision Tree

```
cross-module-git-boundary/
├── module-path/
│   ├── child-path/
│   │   └── child-module-unchanged/          → child path, same git → discover, no warning
│   └── non-child-path/
│       └── git-context/
│           ├── both-null/
│           │   └── both-not-git-discovers/  → neither in git → discover, no warning
│           ├── same-git/
│           │   ├── parent-git-child-same-git-discovers/  → single repo → discover
│           │   └── lifelog-mirror-discovers/             → lifelog layout repro
│           └── different-git/
│               ├── parent-git-child-other-git-warns/     → separate repos → warn + skip
│               └── parent-not-git-child-git-warns/       → parent not git, child git → warn + skip
```

## Test Index

| Leaf | Description |
|------|-------------|
| `module-path/child-path/child-module-unchanged` | Child module path in single git repo: discover nested tests, no warning |
| `module-path/non-child-path/git-context/both-null/both-not-git-discovers` | Neither side in git: discover sibling module tests, no warning |
| `module-path/non-child-path/git-context/same-git/parent-git-child-same-git-discovers` | Single git repo, non-child paths: discover nested module (lifelog case) |
| `module-path/non-child-path/git-context/same-git/lifelog-mirror-discovers` | Lifelog module layout mirror in single git repo: discover `skill-cli` tests |
| `module-path/non-child-path/git-context/different-git/parent-git-child-other-git-warns` | Separate git repos: warn with `different git repository`, skip child |
| `module-path/non-child-path/git-context/different-git/parent-not-git-child-git-warns` | Parent not in git, child in git: warn with `git repository mismatch`, skip child |

## How to Run

```sh
doctest test ./tests/discovery/cross-module-git-boundary/...
doctest test ./tests/test/dotdotdot/cross-go-mod/...
doctest test ./tests/test/dotdotdot/git-boundary/...
```

```go
import (
	"github.com/xhd2015/doctest/libdoc/cli"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type Request struct {
	Args	[]string
	Env	[]string
	WorkDir	string
	Timeout	time.Duration
	Bin	string
	UseCLI	bool
}
type Response struct {
	ExitCode	int
	Stdout		string
	Stderr		string
	Err		error
}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	if !req.UseCLI {
		var stdout, stderr bytes.Buffer
		err := cli.RunWithWriters(&stdout, &stderr, req.Args)
		resp := &Response{Stdout: stdout.String(), Stderr: stderr.String(), Err: err}
		if err != nil {
			resp.ExitCode = 1
			if resp.Stderr == "" {
				resp.Stderr = err.Error()
			}
			return resp, nil
		}
		return resp, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	if req.Timeout <= 0 {
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
	}
	defer cancel()
	bin := req.Bin
	if bin == "" {
		return nil, fmt.Errorf("UseCLI requires req.Bin")
	}
	cmd := exec.CommandContext(ctx, bin, req.Args...)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	}
	if len(req.Env) > 0 {
		cmd.Env = append(os.Environ(), req.Env...)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String(), Err: err}
	if err == nil {
		return resp, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		resp.ExitCode = exitErr.ExitCode()
		return resp, nil
	}
	if ctx.Err() != nil {
		return resp, ctx.Err()
	}
	return resp, err
}

```
