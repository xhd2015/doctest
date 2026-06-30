# Imports Process: Library-based import formatting

## Version
0.0.2


Tests that `WriteGeneratedCase` correctly uses `golang.org/x/tools/imports.Process()` instead of the `goimports` binary, removing the external binary dependency.

## Decision Tree

```
tests/imports-process/
├── DOCTEST.md                      # This file
├── SETUP.md                        # Root: builds doctest binary, run helpers
├── unused-import-removed/          # R1: SETUP.md has unused import → removed by imports.Process
│   ├── SETUP.md                    # Creates test tree with unused import, runs doctest test
│   └── ASSERT.md                   # Assert exit code 0, generated code compiles
├── syntax-error-reported/          # R2: SETUP.md has syntax error → imports.Process error reported
│   ├── SETUP.md                    # Creates test tree with broken Go code
│   └── ASSERT.md                   # Assert error message, no corrupted file written
└── shared-type-named-results-closure/  # R3: helper (port, alt int) → closure needs parens
    ├── SETUP.md                    # Creates tree with shared-type named results
    └── ASSERT.md                   # Assert generation succeeds (exit 0, no format error)
```

## Test Index

| Leaf | Description |
|------|-------------|
| `unused-import-removed` | SETUP.md imports `"fmt"` but doesn't use it → `imports.Process` removes it → test compiles and passes |
| `syntax-error-reported` | SETUP.md contains invalid Go (unclosed string) → `imports.Process` fails → clean error reported |
| `shared-type-named-results-closure` | Root helper `(port, alt int)` must not emit `func(...) port int, alt int` → inner `doctest test` passes |

## How to Run

```sh
doctest test -v ./tests/imports-process/
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
	WorkDir	string
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
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()

	bin := req.Bin
	if bin == "" {
		return nil, fmt.Errorf("req.Bin is not set")
	}
	cmd := exec.CommandContext(ctx, bin, req.Args...)
	cmd.Dir = req.WorkDir
	cmd.Env = append(os.Environ(), req.Env...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
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
