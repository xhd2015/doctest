# Go Test Timeout Message

## Version
0.0.2


Doc-style tests for clear user-facing messaging when nested `go test`
panics with `test timed out after <duration>`, including cancelled-leaf
accounting on the FAIL line and color accents (Error red, hint gray,
`N cancelled` orange).

Does **not** change default timeout policy: only verifies messaging when
the user passes `--timeout` and the suite actually times out, and that
fast-passing suites do not emit a false timeout Error or `cancelled`.

## DSN (Domain Specific Notion)

### Participants

- **`doctest test`** — CLI subcommand that discovers runnable leaves, builds
  generated go-test packages, runs them with optional `-timeout`, and reports
  progress / failures / final PASS|FAIL.
- **Test tree** — temp fixture hierarchy the command targets (multi-leaf sleep
  tree for timeout leaves, or a fast-pass leaf).
- **`go test`** — underlying runner; when the package exceeds `-timeout`, it
  panics with `test timed out after <d>` (often buried under JSON/filtered
  output unless doctest surfaces it). Leaves that never reach pass/fail/skip
  are **cancelled** relative to discovery **planned** count.
- **Stdout / Stderr** — combined user-visible sinks; timeout Error/hint and
  FAIL accounting must be obvious without digging through a full goroutine dump.

### Behaviors

- **Timeout flag** — `doctest test --timeout=2s` forwards `-timeout=2s` to
  `go test` (flag forwarding covered elsewhere; this tree checks messaging).
- **Surface timeout** — when the nested suite times out, exit ≠ 0 and locked
  wording is visible:

  ```
  Error: go test timed out after 2s
  hint: increase with -timeout=DURATION (e.g. -timeout=30m; -timeout=0 disables)
  ```

- **Cancelled accounting** — on timeout with cancelled > 0:

  - **planned** = discovery leaf count (before go_test rewrites Total to actual_run)
  - **cancelled** = `max(0, planned − pass − fail − skipCount)`
  - FAIL denom uses planned: `FAIL (<passed>/<planned>, <N> cancelled)`
  - v1: show `N cancelled` only (omit `t.Skip` phrase on the timeout FAIL line)

- **Progress line** — quiet compact line stays finished-only:
  `(N Run, N Pass, N Fail, N Cached)` — **no** `Cancelled` segment.
- **Print order** — progress → fail dumps → Error/hint → PASS/FAIL
  (Error/hint **before** the final FAIL line on the user-facing stream).
- **Color (when on)** — Error **red**, hint **gray**, `N cancelled` segment
  **orange** (warning accent, e.g. 256-color `38;5;208`); whole `FAIL (…)`
  token stays red with orange nested on the cancelled phrase preferred.
- **No color** — `--no-color` / ColorNever: same wording, no ANSI.
- **No false positive** — a fast-passing tree must not print
  `Error: go test timed out` and must not mention `cancelled`.

## Decision Tree

```
timeout-message
├── surfaces ──── multi-sleep + --timeout=2s --no-color → Error/hint + FAIL (p/P, N cancelled), order  [RED]
├── color ─────── multi-sleep + --timeout=2s --color → Error red, hint gray, cancelled orange         [RED]
└── fast-pass ─── 1-pass leaf, normal timeout → exit 0, no timeout Error, no cancelled               [GREEN]
```

Split factor (most significant first): **timeout fires vs not**. Under timeout,
secondary axis is **color mode** (`--no-color` vs `--color`).

## Test Index

| Leaf | Scenario | Exit | Output signal | Classic |
|------|----------|------|---------------|---------|
| `surfaces` | 3 sleep leaves; `doctest test --timeout=2s --no-color` | ≠ 0 | Locked Error + hint; `FAIL (0/3, N cancelled)` with N>0; progress has no `cancelled`; Error/hint before FAIL | RED until cancelled + order |
| `color` | same multi-sleep; `--timeout=2s --color` | ≠ 0 | Error red, hint gray, `N cancelled` orange; plain wording as surfaces | RED until color helpers |
| `fast-pass` | Temp 1-pass tree; default/generous timeout; `--no-color` | 0 | Must **not** contain `Error: go test timed out` or `cancelled` | expect GREEN |

## How to Run

```sh
doctest vet ./tests/test/timeout-message
doctest test --label heavy ./tests/test/timeout-message/...
# or all labels:
doctest test --label-all ./tests/test/timeout-message/...
```

Classic TDD: expect **RED** on `surfaces` and `color` until implementer lands
planned FAIL denom + `, N cancelled`, timeout Error/hint coloring, orange
cancelled accent, and Error/hint-before-FAIL print order. `fast-pass` should
stay GREEN.

```go
import (
	"github.com/xhd2015/doctest/libdoc/cli"
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
	UseCLI	bool
}
type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
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
