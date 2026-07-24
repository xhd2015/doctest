# End-of-Run Summary Line

## Version
0.0.2


Doc-style tests for the aggregated `PASS(x/y)` / `FAIL(x/y)` / `no tests`
summary printed after `doctest test` completes.

## DSN (Domain Specific Notion)

### Participants

- **`doctest test`** — CLI subcommand that discovers runnable leaves, builds
  generated go-test packages, runs them, and reports progress.
- **Test tree** — temp or fixture directory hierarchy the command targets.
- **Runner** — aggregates stats across multiple directory arguments within one
  invocation.
- **Stdout / Stderr** — output sinks; summary goes to stdout when cases exist,
  `no tests` goes to stderr when none are found.

### Behaviors

- **Discover** — walk arguments (plain dirs, `...`) and count runnable leaves
  (`Total`).
- **Execute** — run each leaf package; count successes including cache hits
  (`Passed`).
- **Progress** — non-verbose mode prints dots then per-suite `(N Run, N Pass, N Fail, N Cached)`.
- **Failure detail** — `FAIL\t...` lines and go-test failure blocks print
  before the final summary when cases fail.
- **Prepare failure (multi-tree)** — when `./...` has a tree that fails prepare
  and a sibling that runs, overall summary is **`FAIL(p/t)`** (never `PASS`
  when prepare failed); process error is labeled **`prepare failed:`**.
- **Summarize** — after all per-dir output, print one line: `PASS(p/t)` or
  `FAIL(p/t)` on stdout; or `no tests` on stderr when `Total == 0`. Any
  non-nil run error (including prepare-only) forces FAIL summary when cases ran.
- **Color** — when enabled, wrap the entire summary token in green (pass) or
  red (fail); `no tests` is never colored.

## Decision Tree

```
summary-line
├── outcome (runnable cases discovered)
│   ├── all-pass ────────────── 3 pass leaves → PASS(3/3) on stdout
│   ├── mixed-pass-fail ─────── 2 pass + 1 fail → FAIL(2/3) after FAIL lines
│   ├── single-pass ─────────── 1 pass leaf → PASS(1/1)
│   ├── verbose-all-pass ────── -v → PASS(1/1) after go test -v output
│   ├── color-pass ──────────── --color → green PASS(...)
│   ├── color-fail ──────────── --color + 1 fail → red FAIL(...)
│   ├── color-disabled ──────── --no-color → plain PASS/FAIL, no ANSI
│   ├── multi-dir-aggregate ─── two dirs (2+1) → single PASS(3/3) at end
│   └── multi-tree-prepare-fail  good+bad under ./... → FAIL + prepare failed:
└── no-cases (zero runnable leaves)
    └── no-cases ────────────── stderr "no tests", no stdout summary
```

## Test Index

| Leaf | Expected summary |
|------|------------------|
| `all-pass` | `PASS(3/3)` on stdout after dots only |
| `mixed-pass-fail` | `FAIL(2/3)` on stdout after `FAIL\t` lines |
| `single-pass` | `PASS(1/1)` |
| `verbose-all-pass` | `PASS(1/1)` after verbose go-test output |
| `no-cases` | `no tests` on stderr, no summary on stdout |
| `color-pass` | green-wrapped `PASS(1/1)` with `--color` |
| `color-fail` | red-wrapped `FAIL(0/1)` with `--color` |
| `color-disabled` | plain `PASS(1/1)` with `--no-color` |
| `multi-dir-aggregate` | single `PASS(3/3)` after two dir runs |
| `multi-tree-prepare-fail` | `FAIL(p/t)` + stderr `prepare failed:` (never PASS) |

## How to Run

```sh
doctest vet ./tests/test/summary-line
doctest test ./tests/test/summary-line
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
	"regexp"
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
