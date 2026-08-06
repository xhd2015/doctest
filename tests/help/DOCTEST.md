# `doctest` help surfaces — in-process CLI

## Version
0.0.2

**Layer model (coverage backfill):**

| Layer | Share | Where |
|-------|-------|--------|
| **In-process CLI** | **all** | every leaf — harness `Run` calls `cli.RunWithWriter` (stdout capture) then `cli.Run`; no product binary, no `testbin` |

Nested root: does **not** inherit workspace binary Run from `tests/DOCTEST.md`.
All leaves are **unlabeled** (fast). Completeness: same four help scenarios as before.

Out of scope: product feature changes; `tests/vet`, `tests/changed`, `tests/skill`.

# DSN (Domain Specific Notion)

### Participants

- **Harness** — invokes `cli.RunWithWriter(w, args)` so user-facing help text is
  captured into a buffer (Parallel-safe; does not reassign `os.Stdout`).
- **CLI dispatcher (`cli.Run`)** — routes top-level and scoped `--help` to the
  same usage strings the product binary prints.
- **Usage texts** — top-level command list; per-command flag docs for `test`,
  `build`, and `agent generate`.

### Behaviors

- **Top-level help** — `doctest --help` lists major subcommands (`agent`, `vet`,
  `design`, `build`, `test`, `skill`, `list`).
- **Test options** — `doctest test --help` documents runner flags (verbose, count,
  timeout, color, cold-cache, Go-style profile/cover/race, `-overlay`/`--overlay`)
  and does not document removed experiment flags.
- **Build options** — `doctest build --help` documents runner options (`-v`,
  `--rm`, `--gen-dir`, `-count`).
- **Agent generate** — `doctest agent generate --help` documents `<idea>`, `-d`,
  `--dir`.

### Pipeline sketch

```
# all leaves (in-process)
req.Args (e.g. ["--help"] | ["test","--help"] | ...)
  -> cli.RunWithWriter(&buf, args)
       -> withTestStdout(buf, cli.Run)
  -> Response{Stdout: buf.String(), ExitCode from err}
```

## Decision Tree

```
tests/help/
├── DOCTEST.md
├── SETUP.md
├── top-level/              doctest --help
├── test-options/           doctest test --help
├── build-options/          doctest build --help
└── agent-generate/         doctest agent generate --help
```

## Test Index

| Leaf | Args | Expected markers |
|------|------|------------------|
| `top-level` | `--help` | `Usage: doctest`, `agent`, `vet`, `design`, `build`, `test`, `skill`, `list` |
| `test-options` | `test --help` | runner flags + profile/cover + `-overlay`/`--overlay`; no removed experiment flags |
| `build-options` | `build --help` | `Usage: doctest build`, `-v`, `--verbose`, `--rm`, `--gen-dir`, `-count` |
| `agent-generate` | `agent generate --help` | `Usage: doctest agent generate`, `<idea>`, `-d`, `--dir` |

## How to Run

```sh
doctest vet ./tests/help/
doctest test ./tests/help/
doctest test ./tests/help/...
```

```go
import (
	"bytes"
	"testing"

	"github.com/xhd2015/doctest/libdoc/cli"
)

// Request drives one help scenario. Leaves set Args only.
type Request struct {
	Args []string // e.g. ["--help"], ["test", "--help"]
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

// Run dispatches help in-process via cli.RunWithWriter (captures cliStdout).
// No testbin, no exec of the product binary.
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	var buf bytes.Buffer
	err := cli.RunWithWriter(&buf, req.Args)
	resp := &Response{
		Stdout: buf.String(),
		Err:    err,
	}
	if err != nil {
		resp.ExitCode = 1
		resp.Stderr = err.Error() + "\n"
		// Mirror process exit: non-zero CLI error is captured, not a harness fail.
		return resp, nil
	}
	return resp, nil
}
```
