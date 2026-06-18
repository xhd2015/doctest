# Scenario

**Feature**: stdin read errors propagate through agent design and implement commands

```
# optional stdin before agent dispatch
cli.Run(args) -> readStdinIfPresent -> runAgentDesign | runAgentImplement

# errors must not be swallowed
broken stdin -> Stat/ReadAll error -> returned to caller
```

## Preconditions
- The `cli` package is importable.
- Tests replace `os.Stdin` with a file/pipe/directory before calling `cli.Run()`.
- The current `readStdinIfPresent()` ignores errors; after the fix it propagates them.

## Steps
1. Child SETUP.md files configure `req.Args` and `req.StdinFile`.
2. Root `Run` replaces `os.Stdin` with the configured file, then calls `cli.Run(req.Args)`.
3. Stdout/stderr are captured to keep test output clean.

## Context
- These tests verify that errors from `os.Stdin.Stat()` and `io.ReadAll()` inside `readStdinIfPresent()` propagate through `cli.Run()` instead of being swallowed.

```go
import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/xhd2015/doctest/libdoc/cli"
)

type Request struct {
	Args      []string
	StdinFile *os.File
}

type Response struct {
	Err    error
	Stdout string
	Stderr string
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
		Err:    cliErr,
		Stdout: stdoutBuf.String(),
		Stderr: stderrBuf.String(),
	}, nil
}
```
