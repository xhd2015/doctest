# `doctest vet` — structure validation (in-process + sparse e2e)

## Version
0.0.2

**Layer model (coverage backfill):**

| Layer | Share | Where |
|-------|-------|--------|
| **L2 doctest in-process** | **mass (≥80%)** | structure, anti-patterns, path/argv patterns — harness `Run` calls `runner.VetArgs` / `validate.Run*` / optional `cli.Run` with fixture trees under `t.TempDir()` |
| **L3 doctest e2e** | **sparse (≤3)** | `e2e/` — real `doctest` binary for help text and verbose process stdout; `label: heavy` |

Default discovery runs L2 only (unlabeled). Use `--label heavy` for binary smokes.

Out of scope: product feature changes; `tests/changed`, `tests/help`, `tests/skill`.

# DSN (Domain Specific Notion)

### Participants

- **Harness** — default path: invokes `runner.VetArgs` (and `validate.Run` /
  `RunWithOptions` underneath) or `cli.Run` in-process against fixture trees.
- **e2e binary** — sparse L3 leaves spawn the product `doctest` binary for
  CLI help wording and verbose stdout that uses process-level `fmt.Printf`.
- **Fixture tree** — temporary directory with `DOCTEST.md` / `SETUP.md` /
  `ASSERT.md` (and optional anti-pattern content) written under `t.TempDir()`.
- **Path resolver** — `./...` / `path/...` discovery via `libdoc/path_resolve`
  (multi-tree and multi-arg wiring through `VetArgs`).
- **Anti-pattern checker** — flags embedded Go string literals, `go test`
  shell-outs, missing SETUP for ASSERT, while skipping `testdata/`.

### Behaviors

- **Valid tree** — minimal legal root → exit 0.
- **Anti-patterns** — embedded Go / go-test shell-out → non-zero + message;
  anti-pattern only under `testdata/` → exit 0 (skipped).
- **Structure** — ASSERT without SETUP → non-zero.
- **Args / paths** — bare `...` rejected; `./...` and `./sub/...` discover
  roots; `./` cwd tree; multi-dir ok; multi-dir with one invalid fails;
  missing dir arg → `vet requires <dir>`.
- **Rename** — top-level `validate` is unknown (`cli.Run` in-process).
- **Help / verbose** — L3 e2e only: usage documents `-v` / `<dir...>` / `./...`;
  `-v` prints `[vet] validating` and file names on process stdout.

### Pipeline sketch

```
# L2 (default)
fixture tree under t.TempDir()
  -> rewrite relative Args against WorkDir when set
  -> runner.VetArgs(flags..., dirs...)   # or cli.Run for non-vet top-level
       -> validate.RunWithOptions(dir, opts)
  -> Response{ExitCode, Stderr from error text}

# L3 (e2e/, label: heavy)
testbin.Ensure -> req.Bin
  -> doctest vet --help | doctest vet -v <dir>
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
│   └── skipped-in-testdata/
├── bare-dot-dot-dot/                   [L2 VetArgs]
├── dot-dot-dot/                        [L2 VetArgs + WorkDir]
├── sub-path-dot-dot-dot/               [L2 VetArgs + WorkDir]
├── dot-slash-dir/                      [L2 VetArgs + WorkDir]
├── multiple-dirs/                      [L2 VetArgs multi]
├── multiple-dirs-one-invalid/          [L2 VetArgs multi]
├── missing-dir/                        [L2 VetArgs]
├── validate-is-unknown/                [L2 cli.Run]
└── e2e/                                [L3 binary, label: heavy]
    ├── help-shows-flags/
    └── verbose-flag/
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
| `bare-dot-dot-dot` | L2 | non-zero; bare `...` message |
| `dot-dot-dot` | L2 | exit 0; discovers sub-a + sub-b |
| `sub-path-dot-dot-dot` | L2 | exit 0; only subp validated |
| `dot-slash-dir` | L2 | exit 0 for `./` |
| `multiple-dirs` | L2 | exit 0 for two valid roots |
| `multiple-dirs-one-invalid` | L2 | non-zero; missing DOCTEST.md |
| `missing-dir` | L2 | non-zero; `vet requires <dir>` |
| `validate-is-unknown` | L2 | non-zero; `unknown command: validate` |
| `e2e/help-shows-flags` | L3 heavy | exit 0; `-v`, `--verbose`, `<dir...>`, `./...` |
| `e2e/verbose-flag` | L3 heavy | exit 0; `[vet] validating` + `SETUP.md` on stdout |

## How to Run

```sh
doctest vet ./tests/vet/
# default discovery: L2 in-process mass (skips label: heavy)
doctest test ./tests/vet/
# sparse binary smokes
doctest test --label heavy ./tests/vet/...
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

func Run(t *testing.T, req *Request) (*Response, error) {
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

	args := append([]string(nil), req.Args...)
	if len(args) == 0 {
		// Treat as bare "vet" with no positional — same as runner.VetArgs(nil).
		if err := runner.VetArgs(nil); err != nil {
			resp.ExitCode = 1
			resp.Stderr = err.Error() + "\n"
			resp.Err = err
		}
		return resp, nil
	}

	// Non-vet top-level commands (e.g. renamed "validate") exercise cli.Run.
	if args[0] != "vet" {
		err := cli.Run(args)
		if err != nil {
			resp.ExitCode = 1
			resp.Stderr = err.Error() + "\n"
			resp.Err = err
		}
		return resp, nil
	}

	// Strip leading "vet"; rewrite relative path args against WorkDir.
	vetArgs := rewriteVetArgs(args[1:], req.WorkDir)
	err := runner.VetArgs(vetArgs)
	if err != nil {
		resp.ExitCode = 1
		resp.Stderr = err.Error() + "\n"
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
