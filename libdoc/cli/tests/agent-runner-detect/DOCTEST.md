# Agent Runner Auto-Detection Tests

## Version
0.0.2


Verify the auto-detection logic in `subagent.Run()` and the removal of hardcoded `"opencode"` defaults in `libdoc/cli/cli.go`.

When `opts.AgentRunner` is empty, the subagent should auto-detect which runner to use. Detection priority:

1. **Env var override**: `DOCTEST_SUBAGENT_AGENT_RUNNER` set → use that value
2. **CODEX_THREAD_ID**: `CODEX_THREAD_ID` set → use `"codex"`
3. **Parent process**: Match parent process name → `"opencode"`, `"pi"`, or `"crush"`
4. **First installed**: `exec.LookPath` in order `"pi"`, `"opencode"`, `"crush"`, `"codex"`
5. **Ultimate fallback**: `"opencode"`

The CLI defaults (`runAgentImplement`, `runAgentDesign`) are changed from `{AgentRunner: "opencode"}` to `{}`, enabling auto-detection.

## Decision Tree

```
tests/agent-runner-detect/                    [Request{Args, Env}, Response{Stdout, Stderr, Err}]
│                                              Run: calls cli.Run(args), captures output and error
├── explicit-runner/                           [--agent-runner=idonotexist → no auto-detection]
│   ├── SETUP.md                               [prepends agent,implement,test,--agent-runner=idonotexist]
│   └── ASSERT.md                              → err "unknown agent runner id: idonotexist"
│
└── auto-detect/                               [no --agent-runner flag → auto-detection enabled]
    ├── SETUP.md                               [prepends agent,implement,test (no --agent-runner)]
    ├── env-override-beats-codex/              DOCTEST_SUBAGENT_AGENT_RUNNER=idonotexist + CODEX_THREAD_ID=abc
    │   ├── SETUP.md                           → err "idonotexist" (env var beats CODEX_THREAD_ID)
    │   └── ASSERT.md
    └── codex-thread-id/                       CODEX_THREAD_ID=abc (no env override)
        ├── SETUP.md                           → err "codex" (CODEX_THREAD_ID triggers codex)
        └── ASSERT.md
```

## Test Index

### explicit-runner — 1 leaf
| Leaf | Description |
|------|-------------|
| `explicit-runner` | `--agent-runner=idonotexist` explicitly passed → error `unknown agent runner id: idonotexist` |

### auto-detect — 2 leaves
| Leaf | Description |
|------|-------------|
| `auto-detect/env-override-beats-codex` | `DOCTEST_SUBAGENT_AGENT_RUNNER=idonotexist` + `CODEX_THREAD_ID=abc` → error `idonotexist` proves env var beats CODEX_THREAD_ID |
| `auto-detect/codex-thread-id` | `CODEX_THREAD_ID=abc` → error mentions `codex` (CODEX_THREAD_ID detection works) |

Total: **3 leaves**.

## How to Run

```sh
doctest test ./libdoc/cli/tests/agent-runner-detect/
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
	// t.Setenv restores env when this test ends (required under workspace suite:
	// os.Setenv(PATH=tmpdir) would poison later trees in the same process).
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
