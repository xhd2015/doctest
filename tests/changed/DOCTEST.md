# `--changed` Flag Integration Tests

## Version
0.0.2

Doc-style tests for the `--changed` flag on `doctest test`, `doctest build`,
and `doctest vet`. When set, only doctest leaves affected by added or modified
files in the working tree are processed. Change detection uses
`git.GetOnDiskChangedFiles` from `github.com/xhd2015/gitops/git`.

## DSN (Domain Specific Notion)

### Participants

- **doctest CLI** — the binary under test; exposes `test`, `build`, and `vet`
  subcommands, each accepting `--changed`.
- **Git repository** — provides the working tree; `git.ShowToplevel` resolves
  the repo root and `git.GetOnDiskChangedFiles` returns staged, unstaged, and
  untracked paths.
- **Fixture test tree** — an ephemeral directory of `DOCTEST.md`, `SETUP.md`,
  and `ASSERT.md` files created inside the git repo; leaves are discovered and
  filtered by the CLI.
- **Changed-file mapper** — maps each changed path to affected leaves: root
  `DOCTEST.md` affects all leaves; group `SETUP.md` affects descendant leaves;
  leaf `ASSERT.md` or `testdata/**` under a leaf affects that leaf only.
  Unrelated non-doctest files under a sibling leaf are ignored.

### Behaviors

- **`--changed` on test/build** — discover leaves, intersect with changed-file
  mapping, run or compile only affected leaves; unchanged siblings are skipped.
- **`--changed` on vet** — validate only changed doctest markdown files; root
  `DOCTEST.md` is not blanket-validated when unchanged.
- **No matching changes** — silent stderr and exit 0; with `-v`, print `doctest: <path> (N tests, --changed: 0 tests)`.
- **Non-git directory** — hard error; `--changed` requires a git repository.
- **Help** — each subcommand documents `--changed` in its usage output.

## Decision Tree

```
--changed
├── help
│   ├── test-help      (stdout includes --changed)
│   ├── build-help     (stdout includes --changed)
│   └── vet-help       (stdout includes --changed)
└── git-context
    ├── not-git-repo
    │   ├── test       (hard error, non-zero exit)
    │   ├── build      (hard error, non-zero exit)
    │   └── vet        (hard error, non-zero exit)
    └── in-git-repo
        ├── test
        │   ├── assert-only           (1 leaf runs)
        │   ├── assert-only-dotdotdot (1 leaf runs via ./...)
        │   ├── assert-only-subpath-dotdotdot (1 leaf runs via ./tests/...)
        │   ├── sibling-stray-untracked (1 leaf runs; sibling stray ignored)
        │   ├── sibling-stray-subpath-dotdotdot (1 leaf runs via ./tests/...; sibling stray ignored)
        │   ├── parent-setup          (2 leaves run)
        │   ├── root-doctest          (2 leaves run)
        │   ├── new-untracked-leaf    (1 new leaf runs)
        │   ├── testdata-change       (1 leaf runs)
        │   └── no-matching-changes
        │       ├── clean-tree        (warning, exit 0)
        │       └── non-doctest-only  (warning, exit 0)
        ├── build
        │   └── assert-only           (1 leaf compiled)
        └── vet
            ├── changed-only          (fails on changed SETUP only)
            └── skip-root             (invalid root skipped, exit 0)
```

## Test Leaf Index

| Leaf | Path |
|------|------|
| test help | `help/test-help/` |
| build help | `help/build-help/` |
| vet help | `help/vet-help/` |
| not-git test | `git-context/not-git-repo/test/` |
| not-git build | `git-context/not-git-repo/build/` |
| not-git vet | `git-context/not-git-repo/vet/` |
| test assert-only | `git-context/in-git-repo/test/assert-only/` |
| test assert-only dotdotdot | `git-context/in-git-repo/test/assert-only-dotdotdot/` |
| test assert-only subpath dotdotdot | `git-context/in-git-repo/test/assert-only-subpath-dotdotdot/` |
| test sibling stray untracked | `git-context/in-git-repo/test/sibling-stray-untracked/` |
| test sibling stray subpath dotdotdot | `git-context/in-git-repo/test/sibling-stray-subpath-dotdotdot/` |
| test parent-setup | `git-context/in-git-repo/test/parent-setup/` |
| test root-doctest | `git-context/in-git-repo/test/root-doctest/` |
| test new-untracked | `git-context/in-git-repo/test/new-untracked-leaf/` |
| test testdata | `git-context/in-git-repo/test/testdata-change/` |
| test clean tree | `git-context/in-git-repo/test/no-matching-changes/clean-tree/` |
| test non-doctest | `git-context/in-git-repo/test/no-matching-changes/non-doctest-only/` |
| build assert-only | `git-context/in-git-repo/build/assert-only/` |
| vet changed-only | `git-context/in-git-repo/vet/changed-only/` |
| vet skip-root | `git-context/in-git-repo/vet/skip-root/` |

## How to Run

```sh
doctest vet ./tests/changed
doctest test ./tests/changed
doctest test ./tests/changed -v
```

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

type Request struct {
	Args    []string
	Env     []string
	WorkDir string
	Timeout time.Duration
	Bin     string
}
type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}
func Run(t *testing.T, req *Request) (*Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()

	bin := req.Bin
	if bin == "" {
		return nil, fmt.Errorf("req.Bin is not set")
	}
	cmd := exec.CommandContext(ctx, bin, req.Args...)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	}
	cmd.Env = append(os.Environ(), req.Env...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
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