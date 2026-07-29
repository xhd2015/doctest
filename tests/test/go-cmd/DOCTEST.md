# `--go-cmd` resolution (auto | xgo | go)

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **CLI (`doctest test`)** — accepts optional `--go-cmd=auto|xgo|go` (default
  **auto** when the flag is omitted). Invalid values fail at parse with a clear
  message. No trap-all / trap flags.
- **GoCmd mode** — the selected policy string after parse: empty/omitted/`auto`,
  forced `xgo`, forced `go`, or invalid.
- **Xgo-mock detector** — walks packages reachable from tree entrypoints
  (Setup / Run / Assert packages and their project-module imports) and reports
  whether `github.com/xhd2015/xgo/runtime/mock` is imported **transitively**.
  False negatives to avoid: Run only calls a project helper that imports mock
  (must follow into project packages, not only markdown text for `mock.Patch`).
- **Resolver** — pure function of (mode, needsXgo) → binary **name** `go` or
  `xgo` (never a PATH shim that renames `go`). Forces ignore detection.
- **PATH lookup** — when the resolved name is `xgo`, the binary must be found on
  PATH (or a test-injected search PATH). Missing → actionable error
  (`xgo not found in PATH` style). Invocation is by name (`xgo test` / `go test`).

**Behaviors**

1. **auto** + no xgo mock reachable → resolve to `go`.
2. **auto** + transitive xgo mock from entrypoints → resolve to `xgo`.
3. **`--go-cmd=xgo`** → always `xgo` even when no mock is detected.
4. **`--go-cmd=go`** → always `go` even when mock is detected (mocks may fail at
   runtime; CLI still honors force-go).
5. **Invalid `--go-cmd=…`** → non-zero, clear error (valid set / invalid value).
6. **Force xgo but binary missing** → non-zero, clear “not found in PATH” error.

**Out of scope (this tree)**

- trap-all / trap package flags
- PATH wrappers that shadow `go` with xgo
- Full e2e suite run with a real xgo binary (optional later; label heavy/e2e)

## Decision Tree

```
go-cmd/                                      [--go-cmd policy + resolve]
├── auto/                                    omitted or --go-cmd=auto
│   ├── no-mock/                             no mock in entry graph → go
│   └── transitive-mock/                     Run→helper→xgo/runtime/mock → xgo
├── force-xgo/                               --go-cmd=xgo
│   ├── available/                           always xgo (no mock needed)
│   └── missing/                             xgo not on search PATH → error
├── force-go/                                --go-cmd=go
│   └── with-mock/                           always go despite mock graph
└── invalid/                                 --go-cmd=<bad>
    └── reject/                              parse/resolve error, clear message
```

Split factor at root children: **go-cmd mode** (auto | force-xgo | force-go |
invalid) — largest product switch. Under auto, split by **mock detection**.
Under force-xgo, split by **binary availability**. force-go covers the
interesting override (mock present). invalid covers bad values.

## Test Index

| Leaf | Scenario (parent design) | Expect today (Classic TDD) |
|------|--------------------------|----------------------------|
| `auto/no-mock` | 1 — auto + no mock → `go` | **RED** until detect+resolve |
| `auto/transitive-mock` | 2 — auto + transitive mock → `xgo` | **RED** until transitive detect |
| `force-xgo/available` | 3 — force xgo always | **RED** until force resolve |
| `force-xgo/missing` | 6 — xgo forced, missing → error | **RED** until PATH check |
| `force-go/with-mock` | 4 — force go despite mock | **RED** until force resolve |
| `invalid/reject` | 5 — invalid value → error | **RED** until flag + validation |

## Implementer surface (library hooks)

Leaves call in-process APIs so a real `xgo` binary is **not** required except
the missing-PATH leaf (fake empty search PATH):

| API | Role |
|-----|------|
| `core.Options.GoCmd` | Parsed mode string: `""`/`"auto"` / `"xgo"` / `"go"` |
| `runner.ParseTestOptions` | Accept `--go-cmd=…` / `--go-cmd …`; reject other values |
| `core.DetectXgoMockUsage(modRoot, entryImportPaths []string) (bool, error)` | Transitive import walk for `github.com/xhd2015/xgo/runtime/mock` under the project module |
| `core.ResolveGoTestCmd(mode string, needsXgo bool) (cmd string, error)` | `auto`/empty → `xgo` if needs else `go`; force `xgo`/`go`; invalid → error |
| `core.EnsureGoTestCmdAvailable(cmd, searchPATH string) error` | Look up `cmd` on `searchPATH` (if empty, process PATH). Missing → clear error |

Product code later wires resolve+ensure into the real `go test` / `xgo test`
exec path; this tree seals the contract without full suite e2e.

## How to Run

```sh
cd external/doctest-master-2026-07-28-1
doctest vet ./tests/test/go-cmd/
doctest test ./tests/test/go-cmd/ --label-all
# Classic TDD: leaves RED until GoCmd parse + detect + resolve + ensure land
```

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/runner"
)

// Request drives parse → detect → resolve → optional PATH ensure for --go-cmd.
// Classic TDD: product helpers named below may be missing until implementer lands them.
type Request struct {
	// ParseOnly: only ParseTestOptions(ParseArgs); used by invalid/* leaves.
	ParseOnly bool
	ParseArgs []string // args after "test" subcommand (no "test" prefix)

	// GoCmdFlag is the mode for ResolveGoTestCmd: "", "auto", "xgo", "go", or garbage.
	// When ParseArgs is set and parse succeeds, Run may copy Options.GoCmd here.
	GoCmdFlag string

	// Detect: when true, call DetectXgoMockUsage(ModRoot, EntryImportPaths).
	// When false, use NeedsXgo as an injected detection result.
	Detect           bool
	NeedsXgo         bool
	ModRoot          string
	EntryImportPaths []string

	// CheckAvailable: after resolve, call EnsureGoTestCmdAvailable(cmd, SearchPATH).
	CheckAvailable bool
	// SearchPATH: if non-empty, Ensure searches only this PATH (fake PATH for missing leaf).
	SearchPATH string
}

// Response exposes parse/detect/resolve outcomes without exec'ing go/xgo test.
type Response struct {
	ParsedGoCmd string // Options.GoCmd after successful parse
	NeedsXgo    bool
	ResolvedCmd string // "go" or "xgo"
	ExitCode    int    // 0 ok; 1 on parse/resolve/ensure error
	ErrMsg      string
	ParseErr    string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}

	if len(req.ParseArgs) > 0 {
		opts, _, err := runner.ParseTestOptions(req.ParseArgs)
		if err != nil {
			resp.ExitCode = 1
			resp.ParseErr = err.Error()
			resp.ErrMsg = err.Error()
			return resp, nil
		}
		// Options.GoCmd is the product field implementer adds ("" or "auto"|"xgo"|"go").
		resp.ParsedGoCmd = opts.GoCmd
		if req.GoCmdFlag == "" {
			req.GoCmdFlag = opts.GoCmd
		}
		if req.ParseOnly {
			return resp, nil
		}
	}

	needs := req.NeedsXgo
	if req.Detect {
		var err error
		needs, err = core.DetectXgoMockUsage(req.ModRoot, req.EntryImportPaths)
		if err != nil {
			resp.ExitCode = 1
			resp.ErrMsg = err.Error()
			return resp, nil
		}
	}
	resp.NeedsXgo = needs

	cmd, err := core.ResolveGoTestCmd(req.GoCmdFlag, needs)
	if err != nil {
		resp.ExitCode = 1
		resp.ErrMsg = err.Error()
		return resp, nil
	}
	resp.ResolvedCmd = cmd

	if req.CheckAvailable {
		if err := core.EnsureGoTestCmdAvailable(cmd, req.SearchPATH); err != nil {
			resp.ExitCode = 1
			resp.ErrMsg = err.Error()
			// Keep ResolvedCmd so asserts can see we chose xgo then failed lookup.
			return resp, nil
		}
	}
	return resp, nil
}
```
