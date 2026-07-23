# Agent Logf: Unified Event Stream Logging

## Version
0.0.2


Verify that all event stream output from sub-agents (traceSession, showStatus)
uses `Logf` for consistent timestamped logging `[2006-01-02T15:04:05]`, while
non-event UI framing (borders, headers) continues without timestamps via
`fmt.Fprintf(os.Stdout, ...)`.

The `logf/` subtree is a standalone root (own `DOCTEST.md`) because it calls
`subagent.Logf` in-process rather than shelling out.

## Decision Tree

```
tests/agent-logf/
├── DOCTEST.md                     # This file
├── SETUP.md                       # Root: Request/Response, Run shells out to doctest binary
│
├── logf/                          # === Standalone root (in-process Logf testing) ===
│   ├── DOCTEST.md
│   └── ...
│
├── trace-session/                 # === Shell out: doctest agent implement --trace ===
│   ├── SETUP.md                   # Setup: sets TEST_GROUP=trace-session
│   ├── no-events-file/            # No events.jsonl → (no events yet) via Logf
│   └── with-events/              # Events exist, session finished → event+done lines via Logf
│
└── show-status/                   # === Shell out: doctest agent implement --status ===
    ├── SETUP.md                   # Setup: sets TEST_GROUP=show-status
    ├── session-not-found/         # No session → stderr error, no timestamp
    ├── no-events/                 # Session found, no events → "No events yet" via Logf
    └── with-events/              # Session found, events exist → event lines via Logf
```

## Test Index

### Logf (standalone root — see `logf/DOCTEST.md`)

### traceSession (2 leaves)
| Leaf | Description |
|------|-------------|
| `trace-session/no-events-file` | No events file: "(no events yet)" + Done via Logf, borders without timestamps |
| `trace-session/with-events` | Events present: event lines via Logf, Done via Logf, borders without timestamps |

### showStatus (3 leaves)
| Leaf | Description |
|------|-------------|
| `show-status/session-not-found` | Session not found: stderr error, no timestamp |
| `show-status/no-events` | No events: header block without timestamps, "No events yet" via Logf |
| `show-status/with-events` | Events present: header block without timestamps, event lines via Logf |

## How to Run

```sh
# All shell-out tests (show-status + trace-session):
doctest test tests/agent-logf/

# In-process Logf tests:
doctest test tests/agent-logf/logf/
```

```go
import (
	"github.com/xhd2015/doctest/libdoc/cli"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

type Request struct {
	Args	[]string
	Env	[]string
	WorkDir	string
	Timeout	time.Duration
	Bin	string
	UseCLI	bool
}
type Response struct {
	ExitCode	int
	Stdout		string
	Stderr		string
	Err		error
}
func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	// Env (e.g. DOCTEST_DEBUG_SESSION_HOME) needs a child process — never Setenv.
	// WorkDir alone can stay L2 but these leaves almost always set Env.
	needProc := req.UseCLI || len(req.Env) > 0 || req.WorkDir != ""
	if !needProc {
		var stdout, stderr bytes.Buffer
		err := cli.RunWithWriters(&stdout, &stderr, req.Args)
		resp := &Response{Stdout: stdout.String(), Stderr: stderr.String(), Err: err}
		if err != nil {
			resp.ExitCode = 1
			if resp.Stderr == "" {
				resp.Stderr = err.Error()
			}
			return resp, nil
		}
		return resp, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	if req.Timeout <= 0 {
		cancel()
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Minute)
	}
	defer cancel()
	bin := req.Bin
	if bin == "" {
		return nil, fmt.Errorf("UseCLI/Env require req.Bin")
	}
	cmd := exec.CommandContext(ctx, bin, req.Args...)
	if req.WorkDir != "" {
		cmd.Dir = req.WorkDir
	}
	if len(req.Env) > 0 {
		cmd.Env = append(os.Environ(), req.Env...)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	resp := &Response{Stdout: stdout.String(), Stderr: stderr.String(), Err: err}
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
