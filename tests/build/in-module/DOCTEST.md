# Internal Import Scan + Temp Compile + Gen-Dir Dump

## Version
0.0.2


Tests for automatic handling of parent-module `internal/` imports: doctest scans
assembled generated Go for `<module>/internal/` import paths, compiles in a
temp directory under the parent module, and optionally copies output to
`--gen-dir` as a review dump.

# DSN (Domain Specific Notion)

**Participants**

- **Doctest CLI** — assembles generated Go from markdown specs, scans import
  paths for parent-module `internal/` usage, and chooses compile strategy.
- **Parent Go module** — project under test (e.g. `example.com/app`) with
  `go.mod`, public packages, and `internal/` packages.
- **Import scanner** — inspects real import paths in assembled source (not
  string/comment substrings) for `<parentModule>/internal/`.
- **Compile temp** — ephemeral `.doctest_run_*` directory under `moduleRoot`;
  holds generated `*_test.go` for `go test`/`go build` using parent `go.mod`.
- **Gen-dir dump** — when `--gen-dir` is set, a copy of generated files for
  user review; not the compile root; never receives a nested `go.mod`.
- **Legacy nested module** — used when no internal imports detected: cache or
  outside gen-dir gets `module testcase` + replace.

**Behaviors**

- Internal import detected → compile in `.doctest_run_*` under module root, no
  nested go.mod, remove temp after run; copy to `--gen-dir` if set.
- No internal import → unchanged legacy nested-module behavior.
- `--gen-dir` dump persists after run (unless `--rm`); compile temp always removed.

## Decision Tree

```
in-module/                                    Root: build doctest binary, temp module helpers
│
├── test/                                     Operation: doctest test
│   ├── internal-import-compiles              internal import, no --gen-dir
│   │                                          → temp compile, exit 0, temp removed
│   ├── internal-import-with-gen-dir-dump     --gen-dir inside module
│   │                                          → temp compile, dump has files, no go.mod
│   └── internal-import-outside-gen-dir-dump  --gen-dir outside module
│                                              → temp compile in-module, dump outside, no go.mod
│
├── build/                                    Operation: doctest build
│   └── no-nested-gomod-in-dump               build + --gen-dir inside module
│                                              → dump has *_test.go, no go.mod
│
├── legacy/                                   No internal import
│   └── nested-module-unchanged               public import, gen-dir outside
│                                              → nested testcase go.mod + replace, exit 0
│
└── lifecycle/                                Temp directory lifecycle
    ├── compile-temp-removed                  internal import + --gen-dir
    │                                          → .doctest_run_* gone after run, dump remains
    └── sigint-leaves-compile-temp            SIGINT during writeCases
                                                → .doctest_run_* removed (interrupt cleanup)
```

## Test Index

| # | Leaf | Description |
|---|------|-------------|
| 1 | `test/internal-import-compiles` | Internal import without `--gen-dir` passes via temp compile; temp removed |
| 2 | `test/internal-import-with-gen-dir-dump` | `--gen-dir` inside module: dump has test file, no nested go.mod |
| 3 | `test/internal-import-outside-gen-dir-dump` | `--gen-dir` outside module: dump copied, no go.mod, temp removed |
| 4 | `build/no-nested-gomod-in-dump` | `doctest build` dumps `*_test.go` without nested go.mod |
| 5 | `legacy/nested-module-unchanged` | Public import + outside gen-dir preserves legacy nested module |
| 6 | `lifecycle/compile-temp-removed` | `.doctest_run_*` removed after test; `--gen-dir` dump persists |
| 7 | `lifecycle/sigint-leaves-compile-temp` | SIGINT during `writeCases` removes `.doctest_run_*` under module root |

## Testdata

| Fixture | Location | Used by |
|---------|----------|---------|
| `internal-module` | `testdata/internal-module/` | Single-leaf internal import |
| 40-leaf module | `lifecycle/sigint-leaves-compile-temp/testdata/` | SIGINT repro (`leaf01`…`leaf40`) |

## How to Run

```sh
# Validate tree structure
doctest vet ./tests/build/in-module

# Run all in-module tests
doctest test ./tests/build/in-module/...

# Run feature tests (expect RED before import-scan implementation)
doctest test ./tests/build/in-module/test/...
doctest test ./tests/build/in-module/build/...
doctest test ./tests/build/in-module/lifecycle/...

# Run legacy regression
doctest test ./tests/build/in-module/legacy/...

# Run a specific leaf
doctest test ./tests/build/in-module/test/internal-import-compiles

# SIGINT compile-temp cleanup
doctest test -v ./tests/build/in-module/lifecycle/sigint-leaves-compile-temp
```

```go
import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	// InterruptDuringWriteCases sends SIGINT after the trigger leaf test file
	// appears in verbose stderr, reproducing Ctrl-C during temp compile generation.
	InterruptDuringWriteCases	bool
	InterruptTriggerLeaf		int
}
type Response struct {
	ExitCode	int
	Stdout		string
	Stderr		string
	Err		error
}
func Run(t *testing.T, req *Request) (*Response, error) {
	if req.InterruptDuringWriteCases {
		return runDoctestInterruptedDuringWriteCases(t, req)
	}
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
