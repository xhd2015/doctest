# `doctest agent with` Subcommand Tests

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **system under test** — behavior covered by this tree.

### Behaviors
- **run** — executes the scenarios in this suite.


Verify the new `doctest agent with --agent-runner=RUNNER [--model=MODEL] <prog> <args...>` subcommand.

This subcommand sets `DOCTEST_SUBAGENT_AGENT_RUNNER` (and optionally `DOCTEST_SUBAGENT_MODEL`) in the child process environment, then executes `<prog>` with inherited stdin/stdout/stderr.

## Decision Tree

```
tests/agent-with/                             [Request{Args, Env}, Response{ExitCode, Stdout, Stderr, Err}]
│                                              Run: subprocess doctest with cmd.Env only (no parent Setenv/stdio hijack)
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
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type Request struct {
	Args    []string
	Env     []string // KEY=VAL for the *child* doctest process only (never parent Setenv)
	Bin     string
	Timeout time.Duration
}
type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error // reconstructed from child stderr when exit != 0 (assert compatibility)
}

// mergeChildEnv builds an isolated child environment: start from os.Environ(),
// then apply overrides (including empty values). Does not mutate the parent process.
func mergeChildEnv(overrides []string) []string {
	m := make(map[string]string, 64)
	for _, e := range os.Environ() {
		k, v, _ := strings.Cut(e, "=")
		m[k] = v
	}
	for _, e := range overrides {
		k, v, ok := strings.Cut(e, "=")
		if !ok {
			continue
		}
		m[k] = v
	}
	out := make([]string, 0, len(m))
	for k, v := range m {
		out = append(out, k+"="+v)
	}
	return out
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	// Subprocess isolation: product code may read env; parent process env stays untouched.
	bin := req.Bin
	if bin == "" {
		return nil, fmt.Errorf("req.Bin is not set")
	}
	timeout := req.Timeout
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, req.Args...)
	cmd.Env = mergeChildEnv(req.Env)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()

	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if cmd.ProcessState != nil {
		resp.ExitCode = cmd.ProcessState.ExitCode()
	}
	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			// CLI prints errors to stderr then exits 1 — surface as resp.Err for asserts.
			msg := strings.TrimSpace(resp.Stderr)
			if msg == "" {
				msg = runErr.Error()
			}
			// Prefer last non-empty line (usage may be multi-line).
			lines := strings.Split(msg, "\n")
			for i := len(lines) - 1; i >= 0; i-- {
				if s := strings.TrimSpace(lines[i]); s != "" {
					msg = s
					break
				}
			}
			resp.Err = errors.New(msg)
			return resp, nil
		}
		if ctx.Err() != nil {
			return resp, ctx.Err()
		}
		return resp, runErr
	}
	return resp, nil
}

```
