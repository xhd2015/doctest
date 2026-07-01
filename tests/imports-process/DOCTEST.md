# Imports Process: Library-based import formatting

## Version
0.0.2


Tests that `WriteGeneratedCase` correctly uses `golang.org/x/tools/imports.Process()` instead of the `goimports` binary, removing the external binary dependency.

## DSN (Domain Specific Notion)

The doctest codegen pipeline lowers a markdown doctest tree into a single
generated `*_test.go` file and compiles it. The participants and behaviors:

- **doctest binary** — the tool under test. Reads a doctest tree, assembles a
  Go test source, formats imports, and runs `go test`.
- **doctest tree** — markdown files (`DOCTEST.md`, `SETUP.md`, `ASSERT.md`)
  whose final go blocks carry `Run`/`Setup`/`Assert` funcs plus helper funcs.
- **AssembleTestSource** — emits a `func TestXxx(t *testing.T)` body that
  declares types/consts/vars, lowers each helper to a **func-literal closure**
  (`name := func(params) results { body }`), wires `Run`/`Setup`/`Assert`, then
  invokes them.
- **writeFuncClosure** — renders a helper's closure signature. Must produce
  valid func-literal syntax: result lists need parentheses when there is more
  than one result or any named result, and named results must be preserved so a
  body that assigns to them compiles.
- **helper emission order** — closures are emitted in some order. Unlike
  top-level funcs, func literals cannot reference a closure declared later, so
  callees must precede callers (topological order).
- **imports.Process** — `golang.org/x/tools/imports.Process` formats the
  generated source and removes unused imports, replacing the external
  `goimports` binary.
- **go test** — compiles and runs the generated test; a codegen bug surfaces as
  a build failure (`undefined: ...`, `build failed`) with a non-zero exit code.

```
doctest tree -> AssembleTestSource -> func-literal closures -> imports.Process -> go test
```

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
├── shared-type-named-results-closure/  # R3: helper (port, alt int) → closure needs parens
│   ├── SETUP.md                    # Creates tree with shared-type named results
│   └── ASSERT.md                   # Assert generation succeeds (exit 0, no format error)
├── named-results-assigned-in-body/ # R4: helper body assigns to named results → names must be preserved
│   ├── SETUP.md                    # Creates tree with helper that assigns to named returns
│   └── ASSERT.md                   # Assert exit 0, no undefined-identifier compile error
└── helper-forward-reference/       # R5: caller defined before callee → closures must be topo-sorted
    ├── SETUP.md                    # Creates tree with caller-then-callee helpers
    └── ASSERT.md                   # Assert exit 0, no undefined-callee compile error
```

## Test Index

| Leaf | Description |
|------|-------------|
| `unused-import-removed` | SETUP.md imports `"fmt"` but doesn't use it → `imports.Process` removes it → test compiles and passes |
| `syntax-error-reported` | SETUP.md contains invalid Go (unclosed string) → `imports.Process` fails → clean error reported |
| `shared-type-named-results-closure` | Root helper `(port, alt int)` must not emit `func(...) port int, alt int` → inner `doctest test` passes |
| `named-results-assigned-in-body` | Root helper body assigns to named results `(mainRepo, wtDir, branch string)` → closure must preserve names → inner `doctest test` passes |
| `helper-forward-reference` | Root helpers `caller` then `callee` (caller calls callee) → closures must be topo-sorted so callee precedes caller → inner `doctest test` passes |

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
