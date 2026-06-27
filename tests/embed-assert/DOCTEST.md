# Embedded Assert Local Module — Integration Tests

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Doctest CLI** — discovers doctest leaves, assembles generated Go, detects
  `github.com/xhd2015/doctest/assert` imports, and chooses compile strategy.
- **Embedded assert source** — single concatenated `assert.go` bytes shipped
  inside the doctest binary (`libdoc/assertmod`).
- **Assert cache** — content-addressed directory under
  `$CACHE/doctest/assert-mod/<md5>/` holding `assert.go` + standalone `go.mod`.
- **Nested testcase module** — legacy path when no parent `internal/` imports:
  generated `module testcase` `go.mod` with `replace` for parent module and assert.
- **Internal-compile temp** — `.doctest_run_*` under parent module root with
  `-modfile` copy of parent `go.mod` plus assert `replace` when needed.
- **Gen-dir dump** — optional review copy; never receives nested `go.mod` on
  internal-compile paths.

**Behaviors**

- Assert import detected → materialize cache (write-once) before code generation.
- No assert import → skip cache materialization; existing go.mod/modfile behavior.
- Legacy nested module + assert → append `replace assert => <cache>` to testcase go.mod.
- Internal compile + assert → temp modfile = parent go.mod + assert replace; `-modfile` on go command.
- Import paths in generated code stay `github.com/xhd2015/doctest/assert` (no rewrite).

## Decision Tree

```
embed-assert/
├── compile-strategy/                         [how Go resolves assert at runtime]
│   ├── nested-module/                        [legacy testcase + nested go.mod]
│   │   ├── replace-in-gomod/                 C1: assert replace in generated go.mod
│   │   ├── assert-output-passes/             C2: assert.Output/Match compile and pass
│   │   ├── import-alias-preserved/           C4: aliased assert import survives assembly
│   │   └── no-assert-no-replace/             C3: no assert import → no assert replace
│   └── internal-compile/                     [temp .doctest_run_* + -modfile]
│       ├── internal-and-assert-modfile/      D1: internal + assert via -modfile
│       ├── no-nested-gomod-in-dump/          D2: --gen-dir dump has no nested go.mod
│       └── internal-only-no-assert-replace/  D3: internal only, no assert replace
├── cache/                                    [assert-mod materialization lifecycle]
│   ├── first-run-materializes/               B1: creates cache dir with assert.go + go.mod
│   ├── second-run-idempotent/                B2: second run does not rewrite cache
│   └── no-import-skips/                      B3: no assert import → no new cache entry
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
| `compile-strategy/nested-module/no-assert-no-replace` | C3 — no assert import, no assert replace in go.mod |
| `compile-strategy/internal-compile/internal-and-assert-modfile` | D1 — internal + assert compiles via `-modfile` |
| `compile-strategy/internal-compile/no-nested-gomod-in-dump` | D2 — gen-dir dump has test files, no go.mod |
| `compile-strategy/internal-compile/internal-only-no-assert-replace` | D3 — internal only, modfile has no assert replace |
| `cache/first-run-materializes` | B1 — first run creates `$CACHE/doctest/assert-mod/<md5>/` |
| `cache/second-run-idempotent` | B2 — second run leaves cache bytes/mtimes unchanged |
| `cache/no-import-skips` | B3 — run without assert import does not create assert-mod entry |
| `operation/build/assert-import-compiles` | E1 — `doctest build` succeeds with assert replace |

## How to Run

```sh
doctest vet ./tests/embed-assert/
doctest test ./tests/embed-assert/          # expect RED before implementation
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
)

type Request struct {
	Args		[]string
	Env		[]string
	WorkDir		string
	Timeout		time.Duration
	Bin		string
	OutsideGenDir	string
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