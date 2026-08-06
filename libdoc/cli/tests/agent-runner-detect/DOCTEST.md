# Agent Runner Auto-Detection Tests

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **system under test** — behavior covered by this tree.

### Behaviors
- **run** — executes the scenarios in this suite.


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
│                                              Run: subprocess doctest with cmd.Env only (no parent Setenv)
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
