# Report Progress CLI Tests

## Version
0.0.2


Test the `report-progress` CLI utility dispatched from the doctest binary.
The utility writes progress updates to a JSONL file so the parent process
can monitor implementation progress.

## DSN (Domain Specific Notion)

### Participants
- **report-progress** — CLI utility that appends JSONL progress entries.
- **parent process** — reads the progress file to monitor implementation.

### Behaviors
- **report** — append a structured progress entry to the configured file.

## How to Run

```sh
doctest test -v ./
```

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type Request struct {
	ProgressFile	string // parent-side path for Assert; also PROGRESS_FILE in child Env
	Args	[]string
	Env	[]string
	Timeout	time.Duration
	Bin	string
}
type Response struct {
	ExitCode	int
	Stdout		string
	Stderr		string
	Err		error
}
func Run(t *testing.T, req *Request) (*Response, error) {
	tmp := t.TempDir()
	rpBin := filepath.Join(tmp, "report-progress")
	if out, err := exec.Command("cp", req.Bin, rpBin).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("copy report-progress: %w\n%s", err, string(out))
	}

	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, rpBin, req.Args...)
	// Do not inherit ambient PROGRESS_FILE / DOCTEST_PROGRESS_FILE. Leaves that
	// need them set it explicitly in req.Env. Otherwise shell/agent pollution
	// breaks errors-without-env-var (and any "unset" scenario).
	cmd.Env = append(filterOutEnvKeys(os.Environ(), "PROGRESS_FILE", "DOCTEST_PROGRESS_FILE"), req.Env...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	resp := &Response{
		Stdout:	stdout.String(),
		Stderr:	stderr.String(),
		Err:	err,
	}
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

// filterOutEnvKeys drops KEY=... entries for the given keys (exact name match).
func filterOutEnvKeys(env []string, keys ...string) []string {
	drop := make(map[string]bool, len(keys))
	for _, k := range keys {
		drop[k] = true
	}
	out := make([]string, 0, len(env))
	for _, e := range env {
		key := e
		if i := strings.IndexByte(e, '='); i >= 0 {
			key = e[:i]
		}
		if drop[key] {
			continue
		}
		out = append(out, e)
	}
	return out
}
```
