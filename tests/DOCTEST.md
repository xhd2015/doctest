# Doctest CLI Integration Tests

## Version
0.0.2


These doc-style tests specify the command-level contract for the doctest CLI.
**Default execution is in-process CLI** (`cli.RunWithWriter`) so the suite mass
stays L2. Leaves that set `UseCLI` invoke a real product binary (true e2e) and
must carry `label: e2e`.

Agent-oriented cases configure fake Codex so no real LLM or network backend is required.

### Default `./...` vs full self-test (labels)

- **Default (discovery)**: `doctest test ./...` — unlabeled in-process mass;
  skips `label: e2e` / `heavy` / other labels.
- **Full integration e2e**: `doctest test --label e2e ./...` — process-boundary
  smokes only (sparse).
- **All leaves**: `doctest test --label-all ./...`.

## DSN (Domain Specific Notion)

### Participants
- **`doctest` CLI surface** — default: in-process via `cli.RunWithWriter` (same
  dispatch as production, no subprocess). Optional `UseCLI` builds one shared
  binary via `testbin.Ensure` and runs it as a subprocess (true e2e).
- **Test tree** — a directory hierarchy of `.md` files (DOCTEST.md, SETUP.md,
  ASSERT.md) that the CLI reads, interprets, and executes. It models a decision
  tree of scenarios.
- **Fake Codex** — a stand-in for the LLM backend used during agent operations
  (`design`, `implement`). It returns predetermined responses so no real
  network or AI model is involved.
- **The file system** — the doctest binary reads and writes files; tests observe
  side effects (generated code, output files, mapping data).

### Behaviors
- **`build`** — compiles source code from a doc-style test directory into an
  executable test binary.
- **`test`** — runs the compiled test binary, reports pass/fail per leaf, and
  propagates exit codes.
- **`agent`** — orchestrates design and implement sub-agents. The design agent
  writes tests; the implement agent writes implementation code. They
  communicate through session state and progress reports.
- **`skill`** — exposes embedded documentation (doc-spec, code-spec) to users.
- **`help`** — prints usage information at top-level and per subcommand.
- **`vet`** — inspects and validates test tree structure.

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

	"github.com/xhd2015/doctest/libdoc/cli"
)

type Request struct {
	Args	[]string
	Env	[]string
	WorkDir	string
	Timeout	time.Duration
	Bin	string
	// UseCLI: when true, spawn req.Bin (true e2e). Default false = in-process CLI.
	// Required when Env or WorkDir must be isolated — never use os.Setenv/Chdir
	// under t.Parallel() leaves.
	UseCLI	bool
	// Parent-side paths for helpers/Assert (also put on child Env). Never process Setenv.
	SessionHome	string
	YieldPQBin	string
	TestQFile	string // yield-pending Assert path; not process env

	// Multi-run harness results (parallel-safe). Filled by leaf Setup helpers
	// (doMultiRun) so Assert never reads package-global state under t.Parallel().
	// Each multi-run uses a unique DOCTEST_SESSION_ID for its warmup/measured pair.
	MRFirst		*Response
	MRSecond	*Response
	MRGenDir	string

	// Cold-cache sandbox paths (parallel-safe). Filled by cold-cache Setup helpers
	// instead of a package-global `st` struct.
	CCCacheHome	string
	CCWarmHome	string
	CCColdHome	string
	CCGenDir	string
	CCTestDir	string
	CCMarker	string
}
type Response struct {
	ExitCode	int
	Stdout		string
	Stderr		string
	Err		error
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	if req.Timeout <= 0 {
		req.Timeout = 2 * time.Minute
	}

	// Env / WorkDir need process isolation. Never os.Setenv or os.Chdir here —
	// leaves run under t.Parallel(). Force the subprocess path.
	if !req.UseCLI && (len(req.Env) > 0 || req.WorkDir != "") {
		return nil, fmt.Errorf("Env/WorkDir require UseCLI (subprocess isolation; no process Setenv/Chdir under Parallel)")
	}

	// Default L2: in-process CLI. No process-global env/cwd mutation.
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

	// L3 e2e: real product binary; Env/Dir only on cmd (child process).
	if req.Bin == "" {
		return nil, fmt.Errorf("UseCLI requires req.Bin")
	}
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, req.Bin, req.Args...)
	cmd.Dir = req.WorkDir
	if len(req.Env) > 0 {
		cmd.Env = append(os.Environ(), req.Env...)
	}
	var stdout, stderr bytes.Buffer
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
