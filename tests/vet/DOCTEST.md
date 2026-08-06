# `doctest vet` — structure validation (in-process L2)

## Version
0.0.2

**Layer model (coverage backfill):**

| Layer | Share | Where |
|-------|-------|--------|
| **L2 doctest in-process** | **mass (≥80%)** | structure, anti-patterns, **vacuous Setup / prose-only SETUP**, path/argv, **L3 layer-share budget** — harness `Run` calls `runner.VetArgs` / `validate.Run*` / optional `cli.Run` with fixture trees under `t.TempDir()` |
| **L3 doctest e2e** | **none in this tree** | help/verbose are L2 in-process CLI (no product binary) |

Default discovery runs L2 only (unlabeled).

Out of scope: product feature changes; `tests/changed`, `tests/help`, `tests/skill`; CLI knobs for max L3 %.

# DSN (Domain Specific Notion)

### Participants

- **Harness** — default path: invokes `runner.VetArgs` (and `validate.Run` /
  `RunWithOptions` underneath) or `cli.Run` in-process against fixture trees.
- **Fixture tree** — temporary directory with `DOCTEST.md` / `SETUP.md` /
  `ASSERT.md` (and optional anti-pattern content or multi-leaf labels) under
  `t.TempDir()`.
- **Path resolver** — `./...` / `path/...` discovery via `libdoc/path_resolve`
  (multi-tree and multi-arg wiring through `VetArgs`).
- **Anti-pattern checker** — flags embedded Go string literals, `go test`
  shell-outs, missing SETUP for ASSERT, Parallel-unsafe harness APIs
  (`os.Setenv`/`Unsetenv`/`Clearenv`, `os.Chdir`, `t.Setenv`, `t.Chdir`,
  `os.Stdout`/`Stderr`/`Stdin` reassignment, `syscall.Setenv`/`Unsetenv`),
  while skipping `testdata/`. Read-only use of stdio (e.g. `Fprint(os.Stdout, …)`)
  is allowed.
- **Setup-policy checker** — non-root `func Setup` that is vacuous (`return nil`
  only, or only `_ = …` blank assigns then `return nil`) hard-fails; message
  tells authors to **remove** the Go **code block** (not "implement the
  behavior"). SETUP.md with `# Scenario` prose and **no** Go block is allowed
  (intermediate and leaf). Setup that does real work is allowed.
- **Layer-share checker** — on **full** tree vet only (not `--changed`): after
  structure + anti-patterns, inventory leaves like `doctest list` (`label: e2e`
  ⇒ L3; else L2). If `leaves >= 10` and `100*L3/leaves > 10` → hard fail
  (`MaxL3Pct=10`, `MinLeaves=10`). Implementer skips share when `opts.ChangedOnly`.

### Behaviors

- **Valid tree** — minimal legal root → exit 0.
- **Anti-patterns** — embedded Go / go-test shell-out / Parallel-unsafe APIs →
  non-zero + `anti-pattern:` message naming the API; anti-pattern only under
  `testdata/` → exit 0 (skipped); stdio read/write without reassignment → exit 0.
- **Vacuous Setup** — non-root body only `return nil`, or only blank assigns then
  `return nil` → non-zero; stderr has `remove` + `code block` and a marker
  (`vacuous` / `anti-pattern` / `return nil` / `blank`); not "implement the behavior".
- **Prose-only SETUP** — intermediate/leaf SETUP with Scenario prose and no Go
  block → exit 0.
- **Real Setup** — non-root Setup with real work (e.g. field assign) → exit 0.
- **Structure** — ASSERT without SETUP → non-zero.
- **Layer share (full vet)** — ≥10 leaves and L3 (e2e) share >10% → non-zero +
  message with path, `L3`, share/pct, and `10%`/`max 10`; ≤10% or <10 leaves →
  exit 0; non-`e2e` labels are L2 for share budget.
- **Args / paths** — bare `...` rejected; `./...` and `./sub/...` discover
  roots; `./` cwd tree; multi-dir ok; multi-dir with one invalid fails;
  missing dir arg → `vet requires <dir>`.
- **Rename** — top-level `validate` is unknown (`cli.Run` in-process).
- **Help / verbose** — L2 in-process: usage documents `-v` / `<dir...>` / `./...`;
  `-v` prints `[vet] validating` and file names on injected stdout.

### Pipeline sketch

```
# L2 (default)
fixture tree under t.TempDir()
  -> rewrite relative Args against WorkDir when set
  -> runner.VetArgs(flags..., dirs...)   # or cli.Run for non-vet top-level
       -> validate.RunWithOptions(dir, opts)
            -> structure + anti-patterns + vacuous Setup policy
            -> (full only) L3 share budget if leaves >= 10
  -> Response{ExitCode, Stderr from error text}
```

## Decision Tree

```
tests/vet/
├── DOCTEST.md
├── SETUP.md
├── valid-tree/                         [L2 validate]
├── valid-no-anti-pattern/              [L2 validate]
├── assert-without-setup/               [L2 validate]
├── anti-pattern/                       [L2 validate]
│   ├── embedded-go/
│   ├── go-test-shellout/
│   ├── skipped-in-testdata/
│   ├── os-setenv/                      Parallel-unsafe os.Setenv
│   ├── os-unsetenv/                    Parallel-unsafe os.Unsetenv
│   ├── os-clearenv/                    Parallel-unsafe os.Clearenv
│   ├── os-chdir/                       Parallel-unsafe os.Chdir
│   ├── t-setenv/                       Parallel-unsafe t.Setenv
│   ├── t-chdir/                        Parallel-unsafe t.Chdir
│   ├── stdio-reassign/                 Parallel-unsafe os.Stdout=…
│   ├── stdio-read-ok/                  positive: Fprint(os.Stdout) OK
│   ├── syscall-setenv/                 Parallel-unsafe syscall.Setenv
│   └── syscall-unsetenv/               Parallel-unsafe syscall.Unsetenv
├── vacuous-setup/                      [L2 validate — Setup body policy]
│   ├── vacuous-return-nil/             non-root Setup = return nil → fail + remove/code block
│   ├── vacuous-blank-assign/           non-root Setup = _=… then return nil → same fail class
│   ├── prose-only-setup-ok/            intermediate SETUP prose, no go block → exit 0
│   └── real-setup-ok/                  Setup does real work → exit 0
├── layer-share/                        [L2 validate — L3 e2e budget]
│   ├── within-budget/                  ≥10 leaves, ≤10% e2e → exit 0
│   ├── over-budget/                    ≥10 leaves, >10% e2e → non-zero
│   ├── tiny-tree-skip/                 <10 leaves, high e2e % → exit 0
│   └── heavy-only-is-l2/               ≥10 all non-e2e labeled (e.g. slow) → exit 0
│   # note: --changed skips share (implementer: opts.ChangedOnly); no leaf
├── bare-dot-dot-dot/                   [L2 VetArgs]
├── dot-dot-dot/                        [L2 VetArgs + WorkDir]
├── sub-path-dot-dot-dot/               [L2 VetArgs + WorkDir]
├── dot-slash-dir/                      [L2 VetArgs + WorkDir]
├── multiple-dirs/                      [L2 VetArgs multi]
├── multiple-dirs-one-invalid/          [L2 VetArgs multi]
├── missing-dir/                        [L2 VetArgs]
├── validate-is-unknown/                [L2 cli.Run]
├── help-shows-flags/                   [L2 cli.Run help]
└── verbose-flag/                       [L2 VetArgs -v]
```

## Test Index

| Leaf | Layer | Expected |
|------|--------|----------|
| `valid-tree` | L2 | exit 0 |
| `valid-no-anti-pattern` | L2 | exit 0 |
| `assert-without-setup` | L2 | non-zero; SETUP missing message |
| `anti-pattern/embedded-go` | L2 | non-zero; embedded Go anti-pattern |
| `anti-pattern/go-test-shellout` | L2 | non-zero; go test shell-out anti-pattern |
| `anti-pattern/skipped-in-testdata` | L2 | exit 0 (testdata skipped) |
| `anti-pattern/os-setenv` | L2 | non-zero; `anti-pattern:` + `os.Setenv` |
| `anti-pattern/os-unsetenv` | L2 | non-zero; `anti-pattern:` + `os.Unsetenv` |
| `anti-pattern/os-clearenv` | L2 | non-zero; `anti-pattern:` + `os.Clearenv` |
| `anti-pattern/os-chdir` | L2 | non-zero; `anti-pattern:` + `os.Chdir` |
| `anti-pattern/t-setenv` | L2 | non-zero; `anti-pattern:` + `t.Setenv` |
| `anti-pattern/t-chdir` | L2 | non-zero; `anti-pattern:` + `t.Chdir` |
| `anti-pattern/stdio-reassign` | L2 | non-zero; `anti-pattern:` + `os.Stdout` |
| `anti-pattern/stdio-read-ok` | L2 | exit 0 (Fprint without reassignment) |
| `anti-pattern/syscall-setenv` | L2 | non-zero; `anti-pattern:` + `syscall.Setenv` |
| `anti-pattern/syscall-unsetenv` | L2 | non-zero; `anti-pattern:` + `syscall.Unsetenv` |
| `vacuous-setup/vacuous-return-nil` | L2 | non-zero; remove + code block; not "implement the behavior" |
| `vacuous-setup/vacuous-blank-assign` | L2 | non-zero; same remove/code-block class as return-nil |
| `vacuous-setup/prose-only-setup-ok` | L2 | exit 0; intermediate SETUP prose, no go block |
| `vacuous-setup/real-setup-ok` | L2 | exit 0; Setup with real work |
| `layer-share/within-budget` | L2 | exit 0; 10 leaves / 1 e2e (10% ≤ max) |
| `layer-share/over-budget` | L2 | non-zero; L3 share message (10 leaves / 2 e2e) |
| `layer-share/tiny-tree-skip` | L2 | exit 0; 3 leaves / 1 e2e (MinLeaves skip) |
| `layer-share/heavy-only-is-l2` | L2 | exit 0; 10 slow leaves, no e2e |
| `bare-dot-dot-dot` | L2 | non-zero; bare `...` message |
| `dot-dot-dot` | L2 | exit 0; discovers sub-a + sub-b |
| `sub-path-dot-dot-dot` | L2 | exit 0; only subp validated |
| `dot-slash-dir` | L2 | exit 0 for `./` |
| `multiple-dirs` | L2 | exit 0 for two valid roots |
| `multiple-dirs-one-invalid` | L2 | non-zero; missing DOCTEST.md |
| `missing-dir` | L2 | non-zero; `vet requires <dir>` |
| `validate-is-unknown` | L2 | non-zero; `unknown command: validate` |
| `help-shows-flags` | L2 | exit 0; `-v`, `--verbose`, `<dir...>`, `./...` |
| `verbose-flag` | L2 | exit 0; `[vet] validating` + `SETUP.md` on stdout |

## How to Run

```sh
doctest vet ./tests/vet/
# default discovery: L2 in-process mass
doctest test ./tests/vet/
# vacuous Setup / prose-only policy only
doctest test ./tests/vet/vacuous-setup/
# layer-share only
doctest test ./tests/vet/layer-share/
# full suite
doctest test --label-all ./tests/vet/...
```

```go
import (
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

	"github.com/xhd2015/doctest/libdoc/cli"
	"github.com/xhd2015/doctest/libdoc/path_resolve"
	"github.com/xhd2015/doctest/libdoc/runner"
	"github.com/xhd2015/doctest/libdoc/validate"
)

// Request drives one vet scenario. Default Mode is in-process (VetArgs / cli.Run).
// Set UseCLI for L3 e2e only.
type Request struct {
	Args    []string // e.g. ["vet", dir] or ["vet", "-v", dir] or ["validate", dir]
	Env     []string // e2e only
	WorkDir string   // fixture cwd; relative Args rewritten against this for L2
	Timeout time.Duration
	Bin     string // e2e only: product binary

	// UseCLI: when true, spawn req.Bin (L3 e2e). Default false = in-process APIs.
	UseCLI bool
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	if req.UseCLI {
		return runE2E(t, req)
	}
	return runInProcess(t, req)
}

func runE2E(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req.Bin == "" {
		return nil, fmt.Errorf("UseCLI requires req.Bin")
	}
	timeout := req.Timeout
	if timeout == 0 {
		timeout = 45 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, req.Bin, req.Args...)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	}
	cmd.Env = append(os.Environ(), req.Env...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()

	resp := &Response{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		Err:      runErr,
		ExitCode: 0,
	}
	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			resp.ExitCode = exitErr.ExitCode()
			return resp, nil
		}
		if ctx.Err() != nil {
			return resp, ctx.Err()
		}
		return resp, runErr
	}
	return resp, nil
}

func runInProcess(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{ExitCode: 0}
	var stdout, stderr bytes.Buffer

	args := append([]string(nil), req.Args...)
	if len(args) == 0 {
		args = []string{"vet"}
	}

	// Non-vet top-level (e.g. "validate") or help → full CLI capture.
	if args[0] != "vet" || len(args) == 1 ||
		(len(args) > 1 && (args[1] == "-h" || args[1] == "--help")) {
		err := cli.RunWithWriters(&stdout, &stderr, args)
		resp.Stdout = stdout.String()
		resp.Stderr = stderr.String()
		if err != nil {
			resp.ExitCode = 1
			if resp.Stderr == "" {
				resp.Stderr = err.Error() + "\n"
			}
			resp.Err = err
		}
		return resp, nil
	}

	// Strip leading "vet"; rewrite relative path args against WorkDir.
	// Capture verbose progress on opts.Stdout (Parallel-safe; no process stdout).
	vetArgs := rewriteVetArgs(args[1:], req.WorkDir)
	err := runner.VetArgsWithWriters(vetArgs, &stdout, &stderr)
	resp.Stdout = stdout.String()
	resp.Stderr = stderr.String()
	if err != nil {
		resp.ExitCode = 1
		if resp.Stderr == "" {
			resp.Stderr = err.Error() + "\n"
		}
		resp.Err = err
	}
	return resp, nil
}

// rewriteVetArgs maps relative path / ./... arguments into absolute forms so
// L2 does not depend on process cwd (parallel-safe; no os.Chdir).
func rewriteVetArgs(args []string, workDir string) []string {
	if len(args) == 0 {
		return args
	}
	out := make([]string, len(args))
	for i, a := range args {
		out[i] = rewriteOneArg(a, workDir)
	}
	return out
}

func rewriteOneArg(arg, workDir string) string {
	if arg == "" || strings.HasPrefix(arg, "-") || arg == "..." {
		return arg
	}
	if filepath.IsAbs(arg) {
		return arg
	}
	if path_resolve.IsDotDotDotPattern(arg) {
		base := path_resolve.ExtractBasePath(arg)
		if workDir != "" {
			if base == "." {
				return filepath.Clean(workDir) + string(filepath.Separator) + "..."
			}
			if !filepath.IsAbs(base) {
				base = filepath.Join(workDir, base)
			}
		} else if base == "." {
			return arg
		} else if !filepath.IsAbs(base) {
			if abs, err := filepath.Abs(base); err == nil {
				base = abs
			}
		}
		return filepath.Clean(base) + string(filepath.Separator) + "..."
	}
	// "./" or "." → WorkDir (or abs cwd)
	if arg == "./" || arg == "." {
		if workDir != "" {
			return workDir
		}
		if abs, err := filepath.Abs("."); err == nil {
			return abs
		}
		return arg
	}
	if workDir != "" {
		return filepath.Join(workDir, arg)
	}
	if abs, err := filepath.Abs(arg); err == nil {
		return abs
	}
	return arg
}

// compile-time touch so validate stays importable for leaf helpers if needed.
var (
	_ = validate.Run
	_ = cli.Run
	_ = runner.VetArgs
)
```
