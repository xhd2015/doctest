# readStdinIfPresent Error Propagation Tests

## Version
0.0.2


Verify that errors from `os.Stdin.Stat()` and `io.ReadAll()` inside
`readStdinIfPresent()` propagate to callers (`runAgentImplement`,
`runAgentDesign`) instead of being silently swallowed.

## DSN (Domain Specific Notion)

### Participants
- **`cli.Run` / `cli.TestExported_RunWithStdin`** — dispatches doctest subcommands from `req.Args` with inject stdin.
- **`readStdinIfPresent`** — reads optional stdin when no positional prompt is given.
- **`runAgentDesign` / `runAgentImplement`** — agent entrypoints that call `readStdinIfPresent`.

### Behaviors
- **stdin inject** — tests pass a controlled file/directory via `TestExported_RunWithStdin` (never reassign process stdin).
- **error propagation** — `Stat`/`ReadAll` failures must surface through `cli.Run`, not be swallowed.

## Decision Tree

```
tests/read-stdin-error/                          [Request{Args, StdinFile}]
│                                                Run: TestExported_RunWithStdin(args, stdin)
├── implement/                                   [prepends "agent","implement" to args]
│   │                                            (via runAgentImplement → readStdinIfPresent)
│   ├── stat-error/                              → err != nil (Stat on closed file propagated)
│   └── read-error/                              → err != nil (ReadAll on directory propagated)
└── design/                                      [prepends "agent","design" to args]
    │                                            (via runAgentDesign → readStdinIfPresent)
    ├── stat-error/                              → err != nil (Stat on closed file propagated)
    └── read-error/                              → err != nil (ReadAll on directory propagated)
```

## Test Index

| Leaf | Description |
|------|-------------|
| `implement/stat-error` | Closed file as stdin during `agent implement` — Stat error propagated |
| `implement/read-error` | Directory as stdin during `agent implement` — ReadAll error propagated |
| `design/stat-error` | Closed file as stdin during `agent design` — Stat error propagated |
| `design/read-error` | Directory as stdin during `agent design` — ReadAll error propagated |

## Coverage Notes

- Happy paths (terminal, pipe with data, empty pipe) are covered by `agent-stdin/`.
- This tree fills the gap for **error propagation** — the two ignored errors in `readStdinIfPresent()`.
- Both callers (`runAgentImplement` and `runAgentDesign`) are tested for both error sources.
- Harness never reassigns `os.Stdin`/`Stdout`/`Stderr` (P5).

## How to Run

```sh
doctest test ./libdoc/cli/tests/read-stdin-error/
```

```go
import (
	"os"
	"testing"
	"github.com/xhd2015/doctest/libdoc/cli"
)

type Request struct {
	Args		[]string
	StdinFile	*os.File
}
type Response struct {
	Err	error
	Stdout	string
	Stderr	string
}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	// Inject stdin via TestExported_RunWithStdin — never reassign os.Stdin/Stdout/Stderr.
	// Error-path leaves fail before agent implement/design produce output, so
	// stdout/stderr capture is unnecessary.
	cliErr := cli.TestExported_RunWithStdin(req.Args, req.StdinFile)
	return &Response{
		Err: cliErr,
	}, nil
}
```
