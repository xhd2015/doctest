# readStdinIfPresent Error Propagation Tests

## Version
0.0.2


Verify that errors from `os.Stdin.Stat()` and `io.ReadAll()` inside
`readStdinIfPresent()` propagate to callers (`runAgentImplement`,
`runAgentDesign`) instead of being silently swallowed.

## DSN (Domain Specific Notion)

### Participants
- **`cli.Run`** — dispatches doctest subcommands from `req.Args`.
- **`readStdinIfPresent`** — reads optional stdin when no positional prompt is given.
- **`runAgentDesign` / `runAgentImplement`** — agent entrypoints that call `readStdinIfPresent`.

### Behaviors
- **stdin replacement** — tests swap `os.Stdin` with a controlled file/pipe/directory.
- **error propagation** — `Stat`/`ReadAll` failures must surface through `cli.Run`, not be swallowed.

## Decision Tree

```
tests/read-stdin-error/                          [Request{Args, StdinFile}]
│                                                Run: replaces os.Stdin, calls cli.Run(args)
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

## How to Run

```sh
doctest test ./libdoc/cli/tests/read-stdin-error/
```

```go
import (
	"bytes"
	"io"
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
func Run(t *testing.T, req *Request) (*Response, error) {
	oldStdin := os.Stdin
	if req.StdinFile != nil {
		os.Stdin = req.StdinFile
		defer func() { os.Stdin = oldStdin }()
	}

	oldStdout := os.Stdout
	rOut, wOut, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	os.Stdout = wOut

	oldStderr := os.Stderr
	rErr, wErr, err := os.Pipe()
	if err != nil {
		wOut.Close()
		rOut.Close()
		return nil, err
	}
	os.Stderr = wErr

	cliErr := cli.Run(req.Args)

	wOut.Close()
	wErr.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	io.Copy(&stdoutBuf, rOut)
	io.Copy(&stderrBuf, rErr)
	rOut.Close()
	rErr.Close()

	os.Stdout = oldStdout
	os.Stderr = oldStderr

	return &Response{
		Err:	cliErr,
		Stdout:	stdoutBuf.String(),
		Stderr:	stderrBuf.String(),
	}, nil
}
```
