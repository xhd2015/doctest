# Logf: Timestamped Logging

## Version
0.0.2


Verify that `subagent.Logf` produces `[2006-01-02T15:04:05]` prefixed output
with correct newline handling, format verbs, and special characters.

## DSN (Domain Specific Notion)

### Participants
- **`subagent.Logf`** — formats a timestamped log line and writes it to stdout.
- **probe subprocess** — tiny `go run` program that calls `Logf`; parent captures
  child stdout only (never reassigns `os.Stdout`).

### Behaviors
- **timestamp prefix** — every line starts with `[YYYY-MM-DDTHH:MM:SS] `.
- **newline normalize** — appends `\n` when the message does not already end with one.

## Decision Tree

```
tests/agent-logf/logf/
├── DOCTEST.md                     # This file
├── SETUP.md                       # Root: Request/Response, Run calls subagent.Logf via subprocess
├── without-trailing-newline/      # Message without \n → newline appended
├── with-trailing-newline/         # Message with \n → no double newline
├── empty-message/                 # Empty string → just timestamp + \n
├── format-verbs/                  # Format string with %s/%d verbs and args
└── special-chars/                 # Multiline message, special characters
```

## Test Index

| Leaf | Description |
|------|-------------|
| `without-trailing-newline` | Message without `\n` gets `\n` appended |
| `with-trailing-newline` | Message with `\n` keeps exactly one `\n` |
| `empty-message` | Empty format string produces timestamp + `\n` |
| `format-verbs` | Format verbs (`%s`, `%d`) resolved correctly |
| `special-chars` | Multiline content, special characters preserved |

## How to Run

```sh
doctest test tests/agent-logf/logf/
```

```go
import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

type Request struct {
	Args	[]string
	Env	[]string
	// ModuleRoot is the doctest module root (request-local) for go run probe.
	ModuleRoot	string
}
type Response struct {
	Stdout string
}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	fmtStr := "default"
	for _, e := range req.Env {
		if strings.HasPrefix(e, "LOGF_FORMAT=") {
			fmtStr = e[len("LOGF_FORMAT="):]
			break
		}
	}

	// Subprocess isolation: never reassign os.Stdout. Product Logf writes via
	// fmt.Print with no inject writer, so capture child stdout only.
	if req.ModuleRoot == "" {
		return nil, fmt.Errorf("req.ModuleRoot is not set")
	}
	agentDir, err := agentProDir(req.ModuleRoot)
	if err != nil {
		return nil, err
	}

	dir := t.TempDir()
	mainPath := filepath.Join(dir, "main.go")
	modPath := filepath.Join(dir, "go.mod")
	mainSrc := `package main

import (
	"fmt"
	"os"

	"github.com/xhd2015/agent-pro/agent/subagent"
)

func main() {
	format := os.Getenv("LOGF_FORMAT")
	if format == "" {
		format = "default"
	}
	args := os.Args[1:]
	if len(args) > 0 {
		ia := make([]any, len(args))
		for i, a := range args {
			ia[i] = a
		}
		subagent.Logf("%s", fmt.Sprintf(format, ia...))
		return
	}
	subagent.Logf("%s", format)
}
`
	if err := os.WriteFile(mainPath, []byte(mainSrc), 0644); err != nil {
		return nil, err
	}
	goMod := "module logfprobe\n\ngo 1.21\n\nrequire github.com/xhd2015/agent-pro v0.0.0\n\nreplace github.com/xhd2015/agent-pro => " +
		filepath.ToSlash(agentDir) + "\n"
	if err := os.WriteFile(modPath, []byte(goMod), 0644); err != nil {
		return nil, err
	}

	// go run requires a tidy module graph/go.sum; without tidy, modern Go
	// fails with "updates to go.mod needed; to update it: go mod tidy".
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = dir
	tidy.Env = append(os.Environ(), "GOWORK=off")
	if out, err := tidy.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("logf probe go mod tidy: %w\n%s", err, out)
	}

	cmdArgs := append([]string{"run", "."}, req.Args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off", "LOGF_FORMAT="+fmtStr)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("logf probe failed: %w\nstderr:\n%s", err, stderr.String())
	}
	return &Response{Stdout: stdout.String()}, nil
}

func agentProDir(moduleRoot string) (string, error) {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "github.com/xhd2015/agent-pro")
	cmd.Dir = moduleRoot
	cmd.Env = append(os.Environ(), "GOWORK=off")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("go list agent-pro: %w", err)
	}
	dir := strings.TrimSpace(string(out))
	if dir == "" {
		return "", fmt.Errorf("empty agent-pro dir from go list")
	}
	return dir, nil
}
```
