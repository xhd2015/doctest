# Suite Summary Runtime t.Skip

## Version
0.0.2


Doc-style tests for the end-of-run `PASS`/`FAIL` summary when suite leaves call
`t.Skip` at runtime (go test `-json` Action `skip`). Classic TDD: leaves that
require `, N t.Skip` and succeeded/actual_run fractioning are **RED** until
implementer lands counting + format changes.

Label-discovery skips (`skipped N labeled …`) are **out of scope** — they stay
on their separate block and must not appear as `, N t.Skip`.

## DSN (Domain Specific Notion)

### Participants

- **`doctest test`** — CLI under test; discovers leaves, runs generated packages
  via `go test` (json-backed stats), prints final `PASS`/`FAIL` summary.
- **Fixture tree** — temp doctest tree with always-pass, always-fail, and/or
  `t.Skip` leaves built per outer leaf.
- **Runtime skip leaf** — suite leaf whose Setup or Assert calls `t.Skip`, so
  go test emits a countable suite-leaf Action `skip`.
- **Pass leaf / fail leaf** — ordinary Assert pass or `t.Fatal` fail leaves.
- **Stdout** — receives progress dots (quiet), optional verbose stream, and the
  final aggregate result line (last non-empty line when cases exist).

### Behaviors

- **Count runtime skips** — suite-leaf Action `skip` increments N (same nesting
  rules as pass/fail: workspace `TestDoctestSuite/<tree>/<leaf>`).
- **Fraction when N > 0** — summary uses **succeeded / actual_run** where
  `actual_run = pass + fail` (skips **excluded** from the denominator). Never
  show planned leaf count in `m` of `PASS (s/m, …)`.
- **Suffix when N > 0** — append `, N t.Skip` inside the parentheses before
  `) in …`.
- **No suffix when N = 0** — keep `PASS (N/N)` / `FAIL (p/t)` unchanged (no
  `t.Skip` text).
- **Exit code** — exit 0 when fail=0 even if N > 0; non-zero when fail > 0.
- **Label skips** — discovery label filter remains a separate
  `skipped N labeled …` block; not folded into `, N t.Skip`.

## Decision Tree

```
summary-tskip
├── one-pass-one-tskip ────────── 1 pass + 1 t.Skip, quiet → PASS (1/1, 1 t.Skip), exit 0  [RED]
├── one-fail-one-tskip ────────── 1 fail + 1 t.Skip, quiet → FAIL (0/1, 1 t.Skip), exit ≠ 0 [RED]
├── all-pass-no-tskip ─────────── 2 pass, no skip → PASS (2/2), no t.Skip text [GREEN]
└── verbose-one-pass-one-tskip ── same as first with -v → PASS (1/1, 1 t.Skip) [RED]
```

Split factor: **runtime skip presence + suite outcome** (pass vs fail vs no-skip
regression). Presentation (`-v` vs quiet) is a secondary axis only for the
one-pass-one-tskip scenario.

## Test Index

| Leaf | Fixture | Flags | Exit | Final summary (target) | Classic |
|------|---------|-------|------|------------------------|---------|
| `one-pass-one-tskip` | `ok` + `skip_me` | `--no-color` | 0 | **`PASS (1/1, 1 t.Skip)`** | RED until implement |
| `one-fail-one-tskip` | `z_fail` + `skip_me` | `--no-color` | ≠ 0 | **`FAIL (0/1, 1 t.Skip)`** | RED until implement |
| `all-pass-no-tskip` | 2 always-pass | `--no-color` | 0 | **`PASS (2/2)`** no `t.Skip` | expect GREEN |
| `verbose-one-pass-one-tskip` | same as first | `-v --no-color` | 0 | **`PASS (1/1, 1 t.Skip)`** | RED until implement |

Must **not** accept wrong fractions such as `PASS (1/2, …)` (planned in
denominator) or bare `PASS (1/1)` when a runtime skip occurred.

## How to Run

```sh
doctest vet ./tests/test/summary-tskip
doctest test --label heavy ./tests/test/summary-tskip/...
```

Classic TDD: expect **RED** on leaves that require `, N t.Skip` until implementer
counts json `skip` and formats the summary. `all-pass-no-tskip` should stay GREEN.

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

func Run(t *testing.T, req *Request) (*Response, error) {
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
