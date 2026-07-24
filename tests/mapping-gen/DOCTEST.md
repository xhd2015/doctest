# Mapping-Gen: Cache-Friendly Test Case Generation

## Version
0.0.2


Tests for the 1-to-1 path mapping generation that mirrors the source directory
structure under a mapping-gen cache root — each doctest leaf becomes its own
Go package so changing one test case only invalidates the Go test cache for
that specific leaf.

## DSN (Domain Specific Notion)

### Participants
- **doctest** — generates per-leaf Go packages under a mapping-gen cache root.
- **source module** — the project `go.mod` that supplies replace directives and dependencies.
- **generated module** — the `testcase` module written into the cache with a replace back to the source tree.

### Behaviors
- **mirror** — copy doctest leaves into cache paths that mirror the source tree.
- **scaffold go.mod** — write a generated `go.mod` with replace directives derived from the source module.
- **tidy** — run `go mod tidy` once per stable source `go.mod` fingerprint.
- **invalidate** — refresh generated `go.mod` when the source `go.mod` or `go.sum` changes.

## Decision Tree

```
mapping-gen/                          Root: builds doctest binary, defines helpers
│
├── test/                             Operation: test (go test)
│   │
│   ├── explicit-gen-dir/             GenDir: --gen-dir specified
│   │   ├── per-leaf-packages/        2 leaves under nested grouping dir
│   │   │                              → verifies per-leaf dir, shared go.mod, no hash
│   │   └── with-pkg-under-test/      2 leaves + Package under test declared
│   │                                  → verifies source files copied to each leaf
│   │
│   ├── auto-gen-dir/                 GenDir: not specified (mapping-gen cache)
│   │   ├── single-leaf-runs/         1 leaf, auto cache dir
│   │   │                              → tests run and pass (exit 0)
│   │   ├── per-leaf-cache-isolation/ 2 leaves, modify one, verify cache isolation
│   │   │                              → unchanged leaf cached, changed leaf rebuilds
│   │   ├── nested-rename-parent-passes/ nested fails, dir renamed, parent passes
│   │   │                              → stale nested cache ignored by parent scope
│   │   └── root-gomod-update-regenerates/ root go.mod replace removed after cache warm
│   │                                      → cached go.mod must drop stale replace
│   │
│   └── error-cases/                  Error paths
│       └── no-test-cases-found/      Empty doctest dir (no ASSERT.md leaves)
│                                      → error: no runnable test cases found
│
└── build/                            Operation: build (compile-only, go build)
    ├── explicit-gen-dir/             GenDir: --gen-dir specified
    │   └── per-leaf-packages/        2 leaves under nested grouping dir
    │                                  → verifies per-leaf dir with compile-only code
    └── auto-gen-dir/                 GenDir: not specified
        └── compiles-successfully/    1 leaf, auto temp dir
                                       → go build succeeds (exit 0)
```

## Test Index

| # | Leaf | Description |
|---|------|-------------|
| 1 | `test/explicit-gen-dir/per-leaf-packages` | Verifies per-leaf directory structure with shared go.mod at project root; no doctest.hash file |
| 2 | `test/explicit-gen-dir/with-pkg-under-test` | Verifies source files copied to each leaf dir with _tc package suffix |
| 3 | `test/auto-gen-dir/single-leaf-runs` | Single leaf test passes using auto mapping-gen cache dir |
| 4 | `test/auto-gen-dir/per-leaf-cache-isolation` | Modifying one leaf's ASSERT.md only invalidates that leaf's cache |
| 5 | `test/auto-gen-dir/nested-rename-parent-passes` | Nested leaf fails, dir renamed + fixed, parent run ignores stale nested cache |
| 6 | `test/auto-gen-dir/root-gomod-update-regenerates` | Root go.mod replace removed after cache warm; cached go.mod must regenerate |
| 7 | `test/error-cases/no-test-cases-found` | Empty doctest tree returns error |
| 8 | `build/explicit-gen-dir/per-leaf-packages` | Build mode generates per-leaf dirs with compile-only stubs |
| 9 | `build/auto-gen-dir/compiles-successfully` | Build mode succeeds with auto temp dir |

## How to Run

```sh
# Run all mapping-gen tests
doctest test ./tests/mapping-gen/...

# Run only test command tests
doctest test ./tests/mapping-gen/test/...

# Run only build command tests
doctest test ./tests/mapping-gen/build/...

# Run a specific leaf
doctest test ./tests/mapping-gen/test/explicit-gen-dir/per-leaf-packages
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
	"strings"
	"testing"
	"time"
)

type Request struct {
	Args	[]string
	Env	[]string
	WorkDir	string
	Timeout	time.Duration
	Bin	string
	// GenDir is the absolute --gen-dir path for this leaf (request-local; no package var).
	GenDir	string
	// Multi-run / multi-phase harness state (request-local; Parallel-safe).
	// Filled by leaf Setup helpers so Assert never reads package-global state.
	MRFirst		*Response
	MRSecond	*Response
	MRGenDir	string
	StaleCachePath	string
}
type Response struct {
	ExitCode	int
	Stdout		string
	Stderr		string
	Err		error
}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
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
