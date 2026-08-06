# Always go test -json Suite Stats (Verbose-Safe Counts)

## Version
0.0.2


Doc-style tests for suite Pass/Fail accounting when nested intentional fails
print `FAIL (` lines into the go-test stream. Verbose (`-v`) presentation must
**not** change Passed counts: counting always comes from `go test -json`
package/test events, never from scanning text for `FAIL (` / `FAIL\t`.

Classic TDD: expect **RED** on the verbose nested case until implementer lands
always-json counting on the verbose path. Quiet nested case should already be
GREEN (json path). Real outer fails still produce `FAIL (p/t)` with p < t.

## DSN (Domain Specific Notion)

### Participants

- **`doctest test`** — CLI under test; discovers leaves, runs generated packages
  via `go test`, aggregates `Passed`/`Total`, prints final `PASS (p/t)` /
  `FAIL (p/t)`.
- **Outer fixture tree** — temp doctest tree the harness invokes (2 leaves for
  nested cases; 1 fail leaf for real-fail).
- **Nested intentional-fail child** — separate temp tree with one forced-fail
  leaf; run by the outer `nested_fail_ok` leaf so nested stdout contains
  `FAIL (0/1)` (and often `FAIL\t…`).
- **Outer leaf `pass_leaf`** — always-pass leaf; contributes one success.
- **Outer leaf `nested_fail_ok`** — runs nested `doctest test` on the child,
  expects non-zero exit, **prints nested stdout** into the outer go-test
  stream, then passes Assert (outer leaf is a success).
- **`go test`** — underlying runner; verbose path currently dumps text and
  (until fixed) scans `FAIL ` / `FAIL\t` prefixes for suite fail counts.
- **Stdout / Stderr** — user-visible sinks; final suite summary is the last
  non-empty stdout line when cases exist.

### Behaviors

- **Always-json counting** — workspace and single-tree runs use `go test -json`
  for Pass/Fail/Run suite accounting even when the user passes `-v` /
  `--verbose`. Verbose only affects presentation (stream more `output` detail).
- **Nested FAIL is not outer fail** — when both outer leaves Assert-pass, final
  summary is **`PASS (2/2)`** even if the stream contains nested `FAIL (0/1)`.
- **Same counts quiet vs verbose** — same outer tree without `-v` also yields
  **`PASS (2/2)`**.
- **Real outer fail** — a true failing outer leaf still yields non-zero exit
  and **`FAIL (p/t)`** with p < t (regression guard).

## Decision Tree

```
verbose-json-stats
├── nested-fail-outer-pass-v ──── outer 2 leaves + nested FAIL ( in -v stream → PASS (2/2)
├── nested-fail-outer-pass-quiet  same tree, no -v → PASS (2/2)
└── real-fail-still-fails ─────── 1 forced outer fail → FAIL (0/1), exit ≠ 0
```

Split factor at root: **what drives suite outcome** (nested intentional fail
leaked into stream vs real outer Assert fail). Presentation mode (`-v` vs quiet)
is a secondary axis only for the nested-outer-pass scenario.

## Test Index

| Leaf | Outer fixture | Flags | Exit | Summary |
|------|---------------|-------|------|---------|
| `nested-fail-outer-pass-v` | `pass_leaf` + `nested_fail_ok` (leaks nested `FAIL (`) | `-v --no-color` | 0 | **`PASS (2/2)`** (RED until always-json) |
| `nested-fail-outer-pass-quiet` | same | `--no-color` | 0 | **`PASS (2/2)`** (expect GREEN on json path) |
| `real-fail-still-fails` | 1 forced-fail leaf | `--no-color` | ≠ 0 | **`FAIL (0/1)`** |

## How to Run

```sh
doctest vet ./tests/test/verbose-json-stats
doctest test --label e2e ./tests/test/verbose-json-stats/...
```

Classic TDD: expect **RED** on `nested-fail-outer-pass-v` until implementer
always uses json counting under `-v`. Other leaves should stay GREEN.

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
