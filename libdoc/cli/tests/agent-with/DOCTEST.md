# `doctest agent with` Subcommand Tests

## Version
0.0.2


Verify the new `doctest agent with --agent-runner=RUNNER [--model=MODEL] <prog> <args...>` subcommand.

This subcommand sets `DOCTEST_SUBAGENT_AGENT_RUNNER` (and optionally `DOCTEST_SUBAGENT_MODEL`) in the child process environment, then executes `<prog>` with inherited stdin/stdout/stderr.

## Decision Tree

```
tests/agent-with/                             [Request{Args, Env}, Response{ExitCode, Stdout, Stderr, Err}]
│                                              Run: replaces os.Stdout/Stderr, calls cli.Run(args)
├── errors/                                    [prepends "agent","with" to Args]
│   ├── missing-agent-runner/                  --agent-runner flag present but no value
│   ├── missing-prog/                          no <prog> positional argument
│   ├── model-without-value/                   --model flag present but no value
│   └── prog-not-found/                        <prog> does not exist in PATH
└── execution/                                 [prepends "agent","with","--agent-runner=opencode" to Args]
    ├── basic/                                 echo hello → stdout "hello\n", exit 0
    ├── with-model/                            --model=gpt-4 → env has DOCTEST_SUBAGENT_MODEL=gpt-4
    ├── with-extra-args/                       extra args forwarded to child
    └── exits-with-code/                       child exits 42 → exit code 42 propagated
```

## Test Index

### Errors — 4 leaves
| Leaf | Description |
|------|-------------|
| `errors/missing-agent-runner` | `--agent-runner` without a value → error `--agent-runner requires a value` |
| `errors/missing-prog` | No `<prog>` argument → error `agent with requires <prog>` |
| `errors/model-without-value` | `--model` without a value → error `--model requires a value` |
| `errors/prog-not-found` | Nonexistent program → error `executable file not found in $PATH` |

### Execution — 4 leaves
| Leaf | Description |
|------|-------------|
| `execution/basic` | Runs `echo hello`, stdout contains `hello`, exit code 0 |
| `execution/with-model` | `--model=gpt-4` set, child sees `DOCTEST_SUBAGENT_MODEL=gpt-4` |
| `execution/with-extra-args` | Extra args after `<prog>` forwarded to child |
| `execution/exits-with-code` | Child exits 42, parent exit code is 42 |

Total: **8 leaves**.

## How to Run

```sh
doctest test ./libdoc/cli/tests/agent-with/
```

```go
import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"github.com/xhd2015/doctest/libdoc/cli"
)

type Request struct {
	Args	[]string
	Env	[]string
}
type Response struct {
	ExitCode	int
	Stdout		string
	Stderr		string
	Err		error
}
func Run(t *testing.T, req *Request) (*Response, error) {
	// t.Setenv restores env when this test ends (workspace multi-tree process safety).
	for _, e := range req.Env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			t.Setenv(parts[0], parts[1])
		}
	}

	oldStdout := os.Stdout
	rOut, wOut, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	os.Stdout = wOut
	defer func() { os.Stdout = oldStdout }()

	oldStderr := os.Stderr
	rErr, wErr, err := os.Pipe()
	if err != nil {
		wOut.Close()
		rOut.Close()
		return nil, err
	}
	os.Stderr = wErr
	defer func() { os.Stderr = oldStderr }()

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

	resp := &Response{
		Stdout:	stdoutBuf.String(),
		Stderr:	stderrBuf.String(),
		Err:	cliErr,
	}
	if cliErr == nil {
		return resp, nil
	}
	var exitErr *exec.ExitError
	if errors.As(cliErr, &exitErr) {
		resp.ExitCode = exitErr.ExitCode()
		return resp, nil
	}
	return resp, nil
}
```
