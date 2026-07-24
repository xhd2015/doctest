# Agent Stdin Input Tests

## Version
0.0.2


Verify that `doctest agent implement` and `doctest agent design` accept
prompts from stdin (pipe/heredoc) when no positional arguments are given.

## Decision Tree

```
tests/agent-stdin/                          [Request{Args, Stdin, StdinSource, ReqContent}]
│                                           Run: replaces os.Stdin, calls cli.Run(args)
├── implement/                              [prepends "agent","implement" to args]
│   ├── no-input/                           → err contains "requires <prompt>"
│   ├── stdin-only/                         → err NOT contains "requires <prompt>"
│   ├── positional-only/                    → err NOT contains "requires <prompt>"
│   ├── positional-wins/                    → err NOT contains "requires <prompt>" (pos wins)
│   ├── stdin-heredoc/                      → err NOT contains "requires <prompt>"
│   ├── stdin-empty-pipe/                   → err contains "requires <prompt>" (empty pipe = no prompt)
│   ├── requirement-only/                   → err NOT contains "requires <prompt>"
│   ├── requirement-positional/             → err NOT contains "requires <prompt>"
│   ├── requirement-stdin/                  → err NOT contains "requires <prompt>"
│   └── requirement-missing/                → err contains "read requirement file"
└── design/                                 [prepends "agent","design" to args]
    ├── no-input/                           → err contains "requires <prompt>"
    ├── stdin-only/                         → err NOT contains "requires <prompt>"
    └── positional-only/                    → err NOT contains "requires <prompt>"
```

## Test Index

| Leaf | Description |
|------|-------------|
| `implement/no-input` | No positional args, stdin is terminal → "requires <prompt>" |
| `implement/stdin-only` | Stdin pipe with content, no positional → prompt from stdin |
| `implement/positional-only` | Positional arg present, no stdin → prompt from args |
| `implement/positional-wins` | Both positional and stdin → positional wins, stdin ignored |
| `implement/stdin-heredoc` | Stdin simulates heredoc with multiline content |
| `implement/stdin-empty-pipe` | Stdin is pipe but empty → "requires <prompt>" |
| `implement/requirement-only` | --requirement file, no prompt → requirement as prompt |
| `implement/requirement-positional` | --requirement + positional → combined prompt |
| `implement/requirement-stdin` | --requirement + stdin → combined prompt |
| `implement/requirement-missing` | --requirement to nonexistent file → file read error |
| `design/no-input` | No positional args, stdin is terminal → "requires <prompt>" |
| `design/stdin-only` | Stdin pipe with content → prompt from stdin |
| `design/positional-only` | Positional arg present → prompt from args |

## How to Run

```sh
doctest test ./libdoc/cli/tests/agent-stdin/
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
	Args        []string
	Stdin       string
	StdinSource string // "pipe" | "devnull" | ""
	ReqContent  string
	Bin         string
	Timeout     time.Duration
}
type Response struct {
	Err    error
	Stdout string
	Stderr string
	ExitCode int
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	// Subprocess isolation: feed stdin via cmd.Stdin; never mutate os.Stdin/Stdout/Stderr.
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
	switch req.StdinSource {
	case "pipe":
		cmd.Stdin = strings.NewReader(req.Stdin)
	case "devnull":
		f, err := os.Open(os.DevNull)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		cmd.Stdin = f
	default:
		cmd.Stdin = nil
	}
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
			msg := strings.TrimSpace(resp.Stderr)
			if msg == "" {
				msg = runErr.Error()
			}
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
