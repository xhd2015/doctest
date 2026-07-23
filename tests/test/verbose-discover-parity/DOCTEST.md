# Verbose Discover Parity + Nested-Tree Parent + Planned Header

## Version
0.0.2


Classic TDD locks for three related prepare/output fixes:

1. **`-v` prepare parity** — verbose must not abort prepare after Light →
   label filter → Hydrate succeeded. A parent tree whose intermediate dir only
   holds a nested `DOCTEST.md` (no SETUP on that intermediate) must prepare and
   run the **same** case set under quiet and `-v`.
2. **Full discover pure nested-tree parents** — intermediate directories with
   **no** `SETUP.md` and no leaf ASSERT under the parent tree (only nested
   DOCTEST children, or empty) must not require a SETUP Go block. If
   `SETUP.md` **exists**, still require Go `Setup`.
3. **Planned header under `-v`** — workspace and single-tree always print a
   planned trees/tests line **before** `go test`, including when Verbose.

Out of scope: agent-logf, interrupt-partial-cached, timeout cancelled UX.

## DSN (Domain Specific Notion)

### Participants

- **`doctest test`** — CLI that discovers leaves (Light → filter → Hydrate),
  materializes generated packages, and runs `go test`. Quiet vs Verbose must
  share the same prepare/run case set; Verbose is presentation + optional
  verbose discover dump, not a second hard-fail discovery gate.
- **Full discover** — `DiscoverTreeCases` / `DiscoverTreeCasesVerbose` walk used
  by build/vet and (today) re-run under `-v`. Must accept pure nested-tree
  parent intermediates without SETUP when the file is absent.
- **Parent tree** — DOCTEST root with its own leaves.
- **Intermediate dir** — path under the parent with **no** SETUP and **no**
  parent-tree ASSERT; only a nested directory that owns its own `DOCTEST.md`
  (or empty). Must not poison parent prepare under `-v` or full discover.
- **Nested DOCTEST root** — inheritance firewall; not part of the parent's case
  set when targeting the parent path alone.
- **Planned header** — user-visible line before `cd … && go test …` announcing
  planned work: single-tree `doctest: <path> (N tests)` or workspace
  `doctest: workspace (N trees, M tests)` (or equivalent `… M tests planned`).

### Behaviors

- **Quiet prepare** — Light discover → filter → Hydrate; print planned
  `(N tests)`; run cases. Intermediate nested-only dirs are irrelevant.
- **Verbose prepare** — same case set as quiet; must **not** fail prepare with
  `intermediate/SETUP.md: must have a Go code block` when intermediate has no
  SETUP and only holds a nested DOCTEST. Planned count still printed under `-v`.
- **Full discover** — missing intermediate SETUP is OK for pure nested parents;
  existing SETUP without Go Setup still errors.
- **Workspace `-v`** — print planned trees+tests before go test command line.

## Decision Tree

```
verbose-discover-parity/
├── prepare-parity/                         Fix 1: quiet vs -v same prepare set
│   ├── quiet-ok/                           quiet parent: exit 0, 1 parent leaf
│   ├── verbose-ok/                         -v parent: exit 0; no intermediate SETUP error
│   └── quiet-verbose-same-count/           both modes: exit 0; same planned count (1)
├── full-discover/                          Fix 2: DiscoverTreeCases pure nested
│   ├── pure-nested-parent-ok/              full discover OK; 1 parent case
│   └── setup-exists-no-go-still-errors/    intermediate SETUP without Go still fails
└── planned-header/                         Fix 3: planned line always under -v
    ├── single-tree-verbose/                -v single tree prints N tests planned
    └── workspace-verbose/                  -v ./... prints trees+tests planned
```

Split factors (most significant first): **prepare mode / discover contract**
(quiet vs verbose vs full API) then **planned header surface** (single vs
workspace). Under full discover, secondary is **SETUP absent vs present-without-Go**.

## Test Index

| Leaf | Scenario | Expected before implement | After implement |
|------|----------|---------------------------|-----------------|
| `prepare-parity/quiet-ok` | Parent+nested intermediate; `doctest test` quiet | **GREEN** (Light path) | **GREEN** |
| `prepare-parity/verbose-ok` | Same fixture; `doctest test -v` | **RED** (full re-discover hard-fail on intermediate) | **GREEN** exit 0; no Go-block error for intermediate |
| `prepare-parity/quiet-verbose-same-count` | Run quiet then `-v`; compare planned count | **RED** (verbose fails / drops) | **GREEN** both exit 0; count=1 both |
| `full-discover/pure-nested-parent-ok` | `DiscoverTreeCases` on parent fixture | **RED** (`must have a Go code block`) | **GREEN** err=nil; 1 case |
| `full-discover/setup-exists-no-go-still-errors` | Intermediate SETUP.md without Go block | **GREEN** (still errors) | **GREEN** still errors |
| `planned-header/single-tree-verbose` | 1-pass; `doctest test -v` | **RED** or weak (count missing on first line) | **GREEN** planned `N tests` visible under `-v` |
| `planned-header/workspace-verbose` | 2 trees; `doctest test -v ./...` | **RED** (workspace omits planned under Verbose) | **GREEN** `(N trees, M tests)` before go test |

## How to Run

```sh
doctest vet ./tests/test/verbose-discover-parity
doctest test --label-all ./tests/test/verbose-discover-parity/...
# Fast API leaves only (unlabeled):
doctest test ./tests/test/verbose-discover-parity/full-discover/...
# Nested CLI locks (heavy):
doctest test --label heavy ./tests/test/verbose-discover-parity/...
```

Classic TDD: expect **RED** on verbose prepare, full-discover pure-nested, and
planned-header leaves until implementer lands the three fixes. `quiet-ok` and
`setup-exists-no-go-still-errors` should stay **GREEN**.

```go
import (
	"github.com/xhd2015/doctest/libdoc/cli"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
)

// Op selects Run strategy:
//   "" / "cli"  — subprocess doctest (req.Args)
//   "dual_cli"  — quiet then verbose CLI (QuietArgs / VerboseArgs → Quiet / Verbose)
//   "discover_full" — core.DiscoverTreeCases(DiscoverRoot)
type Request struct {
	Op           string
	Args         []string
	Env          []string
	WorkDir      string
	Timeout      time.Duration
	Bin          string
	UseCLI	bool
	DiscoverRoot string
	QuietArgs    []string
	VerboseArgs  []string
	// Filled by dual_cli Run for Assert
	Quiet   *Response
	Verbose *Response
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error

	// discover_full
	CaseCount   int
	DiscoverErr string
	Cases       []core.TreeCase
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if !req.UseCLI {
		var stdout bytes.Buffer
		err := cli.RunWithWriter(&stdout, req.Args)
		resp := &Response{Stdout: stdout.String(), Err: err}
		if err != nil {
			resp.ExitCode = 1
			resp.Stderr = err.Error()
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
		return nil, fmt.Errorf("UseCLI requires req.Bin")
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


func runDiscoverFull(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req.DiscoverRoot == "" {
		return nil, fmt.Errorf("DiscoverRoot is not set")
	}
	cases, err := core.DiscoverTreeCases(req.DiscoverRoot)
	resp := &Response{CaseCount: len(cases), Cases: cases}
	if err != nil {
		resp.DiscoverErr = err.Error()
	}
	return resp, nil
}

func runCLI(t *testing.T, req *Request, args []string) (*Response, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()

	bin := req.Bin
	if bin == "" {
		return nil, fmt.Errorf("req.Bin is not set")
	}
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Dir = req.WorkDir
	// Child-only env: never process Setenv; append leaf Env on a copy of environ.
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

func runDualCLI(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	quiet, err := runCLI(t, req, req.QuietArgs)
	if err != nil {
		return quiet, err
	}
	verbose, err := runCLI(t, req, req.VerboseArgs)
	if err != nil {
		// Still attach quiet so Assert can compare partial results.
		req.Quiet = quiet
		return verbose, err
	}
	req.Quiet = quiet
	req.Verbose = verbose
	// Aggregate response mirrors verbose (primary under-test mode) with both attached on req.
	agg := &Response{
		ExitCode: verbose.ExitCode,
		Stdout:   verbose.Stdout,
		Stderr:   verbose.Stderr,
		Err:      verbose.Err,
	}
	return agg, nil
}
```
