# Embedded Assert Local Module — Integration Tests

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Doctest CLI** — discovers doctest leaves, assembles generated Go, and chooses
  compile strategy for external vs internal consumer modules.
- **Embedded assert source** — single concatenated `assert.go` bytes shipped
  inside the doctest binary (`libdoc/assertmod`).
- **Assert cache** — content-addressed directory under
  `$CACHE/doctest/assert-mod/<md5>/` holding `assert.go` + standalone `go.mod`.
- **Nested testcase module** — legacy path when no parent `internal/` imports:
  generated `module testcase` `go.mod` with `replace` for parent module, assert,
  and session (always-on for external modules).
- **Internal-compile temp** — `.doctest_run_*` under parent module root with
  `-modfile` copy of parent `go.mod` plus assert/session replaces (always-on).
- **Gen-dir dump** — optional review copy; never receives nested `go.mod` on
  internal-compile paths.

**Behaviors**

- For **external** consumer modules, generation **always** materializes
  assert-mod (and session-mod) and writes `replace .../assert => <cache>` plus
  session replace into the nested gen `go.mod` — even when author SETUP/ASSERT
  do not import assert (shared gen-root / CondTidy hygiene across multi-tree).
- Doctest self-module still omits assert/session submodule replaces
  (`WriteGoMod` special-case).
- Legacy nested module → assert replace points at cache under outside `--gen-dir`.
- Internal compile → temp modfile = parent go.mod + assert/session replaces;
  `-modfile` on go command when always-on flags apply.
- Import paths in generated code stay `github.com/xhd2015/doctest/assert` (no rewrite).

## Decision Tree

```
embed-assert/
├── compile-strategy/                         [how Go resolves assert at runtime]
│   ├── nested-module/                        [legacy testcase + nested go.mod]
│   │   ├── replace-in-gomod/                 C1: assert replace in generated go.mod
│   │   ├── assert-output-passes/             C2: assert.Output/Match compile and pass
│   │   ├── import-alias-preserved/           C4: aliased assert import survives assembly
│   │   └── no-assert-no-replace/             C3: no author assert → still assert replace
│   └── internal-compile/                     [temp .doctest_run_* + -modfile]
│       ├── internal-and-assert-modfile/      D1: internal + assert via -modfile
│       ├── no-nested-gomod-in-dump/          D2: --gen-dir dump has no nested go.mod
│       └── internal-only-no-assert-replace/  D3: internal only, still -modfile (always-on)
├── cache/                                    [assert-mod materialization lifecycle]
│   ├── first-run-materializes/               B1: creates cache dir with assert.go + go.mod
│   ├── second-run-idempotent/                B2: second run does not rewrite cache
│   └── no-import-skips/                      B3: no author assert → still materialize
└── operation/
    └── build/
        └── assert-import-compiles/           E1: doctest build succeeds with assert
```

## Test Index

| Leaf | Scenario |
|------|----------|
| `compile-strategy/nested-module/replace-in-gomod` | C1 — nested go.mod contains assert replace pointing at cache |
| `compile-strategy/nested-module/assert-output-passes` | C2 — subprocess `doctest test` passes with assert.Output |
| `compile-strategy/nested-module/import-alias-preserved` | C4 — `outputassert` alias preserved in generated test |
| `compile-strategy/nested-module/no-assert-no-replace` | C3 — no author assert import, nested go.mod still has assert replace |
| `compile-strategy/internal-compile/internal-and-assert-modfile` | D1 — internal + assert via `-modfile`; cleans `.doctest.mod` + `.doctest.sum` |
| `compile-strategy/internal-compile/no-nested-gomod-in-dump` | D2 — gen-dir dump has test files, no go.mod |
| `compile-strategy/internal-compile/internal-only-no-assert-replace` | D3 — internal only, still uses `-modfile` (always-on assert+session) |
| `cache/first-run-materializes` | B1 — first run creates `$CACHE/doctest/assert-mod/<md5>/` |
| `cache/second-run-idempotent` | B2 — second run leaves cache bytes/mtimes unchanged |
| `cache/no-import-skips` | B3 — run without author assert import still creates assert-mod |
| `operation/build/assert-import-compiles` | E1 — `doctest build` succeeds with assert replace |

## How to Run

```sh
doctest vet ./tests/embed-assert/
doctest test ./tests/embed-assert/
doctest test ./tests/embed-assert/compile-strategy/...
doctest test ./tests/embed-assert/cache/...
doctest test ./tests/embed-assert/operation/build/...
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

	"github.com/xhd2015/doctest/libdoc/cli"
)

type Request struct {
	Args		[]string
	Env		[]string
	WorkDir		string
	Timeout		time.Duration
	Bin		string
	// UseCLI: true only when child Env is required (isolated DOCTEST_CACHE_HOME).
	// Default false = in-process CLI (Parallel-safe; no os.Setenv/Chdir).
	UseCLI		bool
	// Per-leaf sandbox paths (Parallel-safe). Not package globals — this tree
	// runs leaves under t.Parallel() in one suite package.
	ModuleRoot	string
	TestDir		string
	OutsideGenDir	string
	// CacheHome is the isolated DOCTEST_CACHE_HOME for this leaf (child Env + parent path asserts).
	// Never applied via parent process Setenv.
	CacheHome	string
	// Snapshot state for cache/second-run-idempotent (Parallel-safe, not package globals).
	BeforeAssertMtime	int64
	BeforeAssertDigest	[16]byte
	BeforeGoModMtime	int64
	BeforeGoModDigest	[16]byte
}
type Response struct {
	ExitCode	int
	Stdout		string
	Stderr		string
	Err		error
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	if req.Timeout <= 0 {
		req.Timeout = 120 * time.Second
	}
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
	if req.Bin == "" {
		return nil, fmt.Errorf("UseCLI/WorkDir/Env require req.Bin (root Setup should Ensure; Parallel-safe)")
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
