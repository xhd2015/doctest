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
	"testing"
	"time"
	libdocbuild "github.com/xhd2015/doctest/libdoc/build"
)

type Request struct {
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
	cmd.Env = append(os.Environ(), req.Env...)

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
```
