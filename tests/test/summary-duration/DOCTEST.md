# Summary Duration

## Version
0.0.2


Doc-style tests for elapsed-time in both doctest test summaries:
per-suite dot progress `(N Run, N Pass, N Fail, N Cached) in DURATION` and
final aggregate `PASS (p/t) in DURATION` / `FAIL (p/t) in DURATION`.

## DSN (Domain Specific Notion)

### Participants

- **`doctest test`** — CLI subcommand that discovers runnable leaves, builds
  generated go-test packages, runs them, and reports progress.
- **Test tree** — temp directory hierarchy the command targets.
- **Runner** — aggregates stats and wall-clock time across multiple directory
  arguments within one invocation.
- **Per-suite timer** — measures wall time of a single `go test` subprocess
  (non-verbose dot path only for inline summary).
- **Stdout** — receives dot progress, inline per-suite summaries, and the
  final aggregate result line.

### Behaviors

- **Discover** — walk arguments and count runnable leaves (`Total`).
- **Execute** — run each leaf package via `go test`; emit dots in non-verbose
  mode.
- **Inline summarize** — after dots, print `(N Run, N Pass, N Fail, N Cached)
  in DURATION` where `DURATION` is display-formatted wall time for that `go test`
  subprocess (integer sub-second units; ≥1s with at most 2 decimal digits).
- **Aggregate** — runner sums pass/fail counts and measures total invocation
  wall time from start until just before the final summary.
- **Final summarize** — print one line `PASS (p/t) in DURATION` or
  `FAIL (p/t) in DURATION` on stdout.
- **Color** — when enabled, gray-wrap `DURATION` after the closing paren in
  the inline summary; green/red-wrap only the `PASS (p/t)` / `FAIL (p/t)` token on
  the final line; the ` in DURATION` suffix stays plain.

## Decision Tree

```
summary-duration
├── single-pass ─────────────── 1 leaf, non-verbose → inline + final duration
├── all-pass ────────────────── 3 leaves → inline duration after dots
├── verbose-pass ─────────────── -v → no inline summary; final has duration
├── color-pass ───────────────── --color → gray inline duration; green PASS token
├── color-disabled ───────────── --no-color → plain durations, no ANSI
├── multi-dir-aggregate ──────── two dirs → per-dir inline + one final duration
└── slow-leaf ────────────────── ~1s sleep → durations ≥ 1s
```

## Test Index

| Leaf | Mode | Expected |
|------|------|----------|
| `single-pass` | default | `(1 Run, 1 Pass, 0 Fail, 0 Cached) in <dur>` + `PASS (1/1) in <dur>` |
| `all-pass` | default | 3 dots; `(3 Run, 3 Pass, 0 Fail, 0 Cached) in <dur>` + `PASS (3/3) in <dur>` |
| `verbose-pass` | `-v` | No inline summary; `PASS (1/1) in <dur>` after go test -v output |
| `color-pass` | `--color` | Gray inline duration; green `PASS (1/1)`; plain ` in <dur>` |
| `color-disabled` | `--no-color` | Plain durations throughout stdout |
| `multi-dir-aggregate` | two dirs | Per-dir inline summaries with duration; one `PASS (3/3) in <dur>` |
| `slow-leaf` | default | Parsed durations ≥ 1s |

## How to Run

```sh
doctest vet ./tests/test/summary-duration
doctest test ./tests/test/summary-duration
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
	"regexp"
	"strings"
	"testing"
	"time"
	libdocbuild "github.com/xhd2015/doctest/libdoc/build"
)

type Request struct {
	Args	[]string
	Env	[]string
	WorkDir	string
	Timeout	time.Duration
	Bin	string
}
type Response struct {
	ExitCode	int
	Stdout		string
	Stderr		string
	Err		error
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
		Stdout:	stdout.String(),
		Stderr:	stderr.String(),
		Err:	err,
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
