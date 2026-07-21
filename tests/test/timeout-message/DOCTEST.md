# Go Test Timeout Message

## Version
0.0.2


Doc-style tests for clear user-facing messaging when nested `go test`
panics with `test timed out after <duration>`.

Does **not** change default timeout policy: only verifies messaging when
the user passes `--timeout` and the suite actually times out, and that
fast-passing suites do not emit a false timeout Error.

## DSN (Domain Specific Notion)

### Participants

- **`doctest test`** — CLI subcommand that discovers runnable leaves, builds
  generated go-test packages, runs them with optional `-timeout`, and reports
  progress / failures.
- **Test tree** — temp fixture hierarchy the command targets (sleep leaf or
  fast-pass leaf).
- **`go test`** — underlying runner; when the package exceeds `-timeout`, it
  panics with `test timed out after <d>` (often buried under JSON/filtered
  output unless doctest surfaces it).
- **Stdout / Stderr** — combined user-visible sinks; timeout must be obvious
  without digging through a full goroutine dump.

### Behaviors

- **Timeout flag** — `doctest test --timeout=2s` forwards `-timeout=2s` to
  `go test` (flag forwarding covered elsewhere; this tree checks messaging).
- **Surface timeout** — when the nested suite times out, exit ≠ 0 and a clear
  line is visible, e.g. `Error: go test timed out after 2s` (or an equivalent
  `test timed out after 2s` / `timed out after 2s` signal on the fail path).
- **No false positive** — a fast-passing tree must not print
  `Error: go test timed out`.

## Decision Tree

```
timeout-message
├── surfaces ──── sleep ≥3s leaf + --timeout=2s → non-zero + timeout Error visible
└── fast-pass ─── 1-pass leaf, normal/large timeout → exit 0, no timeout Error
```

## Test Index

| Leaf | Scenario | Exit | Output signal |
|------|----------|------|----------------|
| `surfaces` | Temp tree sleeps in Run; `doctest test --timeout=2s --no-color` | ≠ 0 | `Error: go test timed out after …` (or `test timed out after` / `timed out after`) |
| `fast-pass` | Temp 1-pass tree; default or generous timeout | 0 | Must **not** contain `Error: go test timed out` |

## How to Run

```sh
doctest vet ./tests/test/timeout-message
doctest test --label heavy ./tests/test/timeout-message/...
```

Classic TDD: expect **RED** on `surfaces` until implementer surfaces timeout
messages on the JSON/verbose fail path. `fast-pass` should stay GREEN.

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
	cmd.Dir = req.WorkDir
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
