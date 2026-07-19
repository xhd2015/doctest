# Harness migration: `d *session.Doctest` context (Plan P3)

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Author harness** — SETUP / Run / Assert (and package helpers) written in markdown Go
  blocks. After P3, paths and session identity come only from `d *session.Doctest`
  (`d.DOCTEST_ROOT`, `d.DOCTEST_CASE`, `d.DOCTEST_SESSION_ID`), never free package vars.
- **Generated test** — assembler (P2) injects `d := &session.Doctest{...}` and passes it
  into Setup / Run / Assert. Free `DOCTEST_ROOT` / `DOCTEST_SESSION_ID` package vars are gone.
- **Doctest CLI** — builds and runs a fixture tree; this outer suite invokes a freshly
  built binary against temp fixtures that already use the new harness style.
- **Leaf case dir** — absolute path in `d.DOCTEST_CASE`; leaf-local fixtures (e.g.
  `fixture.txt` next to `ASSERT.md`) are read via `filepath.Join(d.DOCTEST_CASE, name)`.
- **Package helper** — shared functions that need paths take `d *session.Doctest` as a
  parameter (no closed-over free vars).

**Behaviors**

- Outer harness of this tree declares `Setup` / `Run` / `Assert` with `d *session.Doctest`
  and resolves the module root as `filepath.Join(d.DOCTEST_ROOT, "..", "..")`.
- A minimal fixture tree whose Setup records `d.DOCTEST_ROOT` (and Assert checks fields)
  compiles and PASSes under `doctest test` (P2 generate already injects `d`).
- A fixture that reads a file next to the leaf via `d.DOCTEST_CASE` PASSes without
  relying on process cwd or free path vars.
- A fixture package helper `joinCase(d, name)` used from Setup PASSes and returns the
  absolute leaf-local path.

```
author harness (Setup/Run/Assert/helpers)
  -> receive d *session.Doctest from generated test
  -> use d.DOCTEST_ROOT / d.DOCTEST_CASE / d.DOCTEST_SESSION_ID
  -> no free DOCTEST_* package vars; no bare relative leaf I/O

outer leaf
  -> write temp fixture tree (author style uses d)
  -> testbin.Ensure(moduleRoot from d.DOCTEST_ROOT)
  -> doctest test <fixture> -> PASS
```

## Decision Tree

```
harness-d-context/                         # nested root; outer harness uses d
├── e2e-fixture-with-d/                    Setup/Run/Assert take d; root path via d
├── leaf-local-via-case/                   leaf fixture.txt via d.DOCTEST_CASE
└── helper-takes-d/                        package helper joinCase(d, name)
```

## Test Index

| Leaf | Description |
|------|-------------|
| `e2e-fixture-with-d` | Temp tree Setup/Run/Assert take `d`; Setup records `d.DOCTEST_ROOT`; `doctest test` exits 0 |
| `leaf-local-via-case` | Temp leaf has `fixture.txt`; Run reads `filepath.Join(d.DOCTEST_CASE, "fixture.txt")`; pass |
| `helper-takes-d` | Fixture defines `joinCase(d *session.Doctest, name string)`; Setup uses it; Assert path matches |

## How to Run

```sh
doctest vet ./tests/harness-d-context/
doctest test -v ./tests/harness-d-context/
# P2 already injects d → fixture leaves may be GREEN without further product work.
# Mass-migrate of other in-repo harnesses is implementer scope (not this tree).
```

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"
)

// Request configures a subprocess run of the doctest CLI against a fixture tree.
type Request struct {
	Args    []string
	Env     []string
	WorkDir string
	Timeout time.Duration
	Bin     string

	// FixtureDir is the absolute path of the temp doctest tree under test.
	FixtureDir string
}

// Response captures CLI process outcomes.
type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()

	bin := req.Bin
	if bin == "" {
		return nil, fmt.Errorf("req.Bin is not set")
	}
	if len(req.Args) == 0 {
		return nil, fmt.Errorf("req.Args is empty")
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
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
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
