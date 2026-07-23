# `doctest skills update` CLI Integration Tests

End-to-end tests for the doctest binary `skills` subcommand and `skills update`
wiring through `libdoc/spec` into `install.HandleUpdateMany`.

# DSN (Domain Specific Notion)

Participants:

- **doctest CLI** — built from `cmd/doctest`; exposes `skills` with `update` and
  per-skill `--install` via the existing `skill` command shape.
- **Skill registry** — `libdoc/spec` maps short names (`tdd`, `designer`, …) to
  on-disk `SkillDirName` values and embedded markdown content.
- **Install/update library** — `github.com/xhd2015/skills/install` performs
  filesystem work; update only touches dirs that already contain `SKILL.md`.
- **Ephemeral project directory** — subprocess cwd for CLI invocations; may
  receive pre-install via `doctest skill <name> --install`.

Behaviors:

- **`doctest skills update`** — walks every registry skill in stable CLI-name
  order; updates installed targets via `InstallTo`; prints `skill not installed:
  <name>` when no target has `SKILL.md` at the resolved scope.
- **`doctest skill tdd --install`** — seeds one skill before batch update tests.
- **`doctest skills --help`** — documents the `update` subcommand.

## Decision Tree

```text
skills-update/
├── run-outcome/
│   ├── none-installed/                  # no installs → not-installed line per skill
│   ├── global-only-needs-global-flag/   # global install, default scope → all not-installed
│   ├── global-up-to-date-reports/       # global install + --global → tdd + others not-installed
│   └── one-skill-installed/             # local tdd → up-to-date + not-installed for others
└── help/
    └── skills-mentions-update/
```

## Test Index

| Leaf | Description |
|------|-------------|
| `run-outcome/none-installed` | No skills installed → exit 0, not-installed line for every registry skill |
| `run-outcome/global-only-needs-global-flag` | Global `tdd` only; default-scope update → not-installed for all skills (incl. `tdd`) |
| `run-outcome/global-up-to-date-reports` | Global `tdd`; `skills update --global` → up-to-date for `tdd` + not-installed for others |
| `run-outcome/one-skill-installed` | Local `tdd` install → up-to-date for `tdd` + not-installed lines for other registry skills |
| `help/skills-mentions-update` | `doctest skills --help` mentions `update` |

## How to Run

```sh
doctest vet ./tests/skills-update
doctest test -v ./tests/skills-update
```

## Version

0.0.2

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
	"testing"
	"time"
)

type PreInstallCLI struct {
	Args []string // e.g. []string{"skill", "tdd", "--install"}
}

type Request struct {
	Args        []string // argv after binary name, e.g. []string{"skills", "update"}
	PreInstalls []PreInstallCLI
	WorkDir     string // set by Setup; subprocess cwd
	// Home is isolated $HOME for global skill installs; passed via child Env only
	// (never t.Setenv — Parallel-incompatible). Asserts use this path too.
	Home string
	Env  []string // extra KEY=VAL for subprocess (e.g. HOME=...)

	Timeout time.Duration
	Bin     string
	UseCLI	bool
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
	WorkDir  string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if !req.UseCLI {
		var stdout bytes.Buffer
		err := cli.RunWithWriter(&stdout, req.Args)
		resp := &Response{Stdout: stdout.String(), Err: err}
		if err != nil {
			resp.ExitCode = 1
			resp.Stderr = err.Error()
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


func childEnv(extra []string) []string {
	return append(append([]string(nil), os.Environ()...), extra...)
}

func runCLI(t *testing.T, bin, dir string, timeout time.Duration, args, env []string) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	cmd.Env = childEnv(env)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w\n%s", err, out)
	}
	return nil
}

func runCLICapture(t *testing.T, bin, dir string, timeout time.Duration, args, env []string) (*Response, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = dir
	cmd.Env = childEnv(env)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	resp := &Response{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		WorkDir:  dir,
		ExitCode: 0,
	}
	if cmd.ProcessState != nil {
		resp.ExitCode = cmd.ProcessState.ExitCode()
	}
	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			return resp, nil
		}
		return resp, runErr
	}
	return resp, nil
}
```