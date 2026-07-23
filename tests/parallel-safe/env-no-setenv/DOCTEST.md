# Parallel-safe: no process env writes (P1)

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Doctest CLI / runner** — one CLI invocation for `doctest test` (and related go
  tool execs). Holds session id and optional cold `GOCACHE` in **options** for
  the run; must not mutate the process environment with `os.Setenv` /
  `syscall.Setenv` / `os.Unsetenv` / `os.Clearenv` for those keys.
- **Product sources** — primarily `libdoc/runner` (and other non-test `libdoc`
  packages). Process env mutation for `DOCTEST_SESSION_ID` and `GOCACHE` is an
  anti-pattern that races under parallel suite leaves.
- **Child `go test` / go tool** — receives `DOCTEST_SESSION_ID` and optional
  isolated `GOCACHE` only via **`cmd.Env`** (key **replace**, not blind append).
- **Nested suite leaf** — generated harness reads `DOCTEST_SESSION_ID` once via
  `syscall.Getenv` into `d.DOCTEST_SESSION_ID` (reads OK; inherit from child env).
- **Tiny fixture tree** — minimal one-leaf doctest used for functional exit-0
  locks under normal session and `--cold-cache`.

**Behaviors**

1. **Static contract (P1 RED until implementer)**: non-test product sources under
   `libdoc/` must not call process env writers for `DOCTEST_SESSION_ID` /
   `DoctestSessionIDEnv` or `GOCACHE`.
2. **Session still works**: `doctest test` on a tiny fixture exits 0; nested
   leaf sees a non-empty session id via suite inject (`d.DOCTEST_SESSION_ID`).
3. **Cold-cache still works**: `doctest test --cold-cache` on a tiny fixture
   exits 0 and announces isolated GOCACHE without requiring process-level
   `Setenv` after the product fix (functional lock stays green across the
   refactor when implementer plumbs `opts.GoCache` → `cmd.Env`).

```
product libdoc sources
  -> must NOT os.Setenv / syscall.Setenv DOCTEST_SESSION_ID | GOCACHE

doctest test [ --cold-cache ] <tiny-fixture>
  -> opts.SessionID / opts.GoCache held once per CLI run
  -> child go test: cmd.Env key-replace (SESSION_ID, GOCACHE)
  -> leaf: syscall.Getenv -> d.DOCTEST_SESSION_ID (read-only)
  -> exit 0
```

## Decision Tree

```
env-no-setenv/                                      [P1 process env anti-pattern]
├── static-source/                                  product source contract
│   └── no-session-gocache-setenv/                  S1: no Setenv SESSION/GOCACHE in libdoc (non-test)
└── functional/                                     behavior locks (session + cold-cache)
    ├── session-tiny-pass/                          F1: tiny fixture exit 0; nested session non-empty
    └── cold-cache-tiny-pass/                       F2: --cold-cache tiny fixture exit 0 + announce
```

## Test Index

| Leaf | Scenario | Expected before implement | After P1 product fix |
|------|----------|---------------------------|----------------------|
| `static-source/no-session-gocache-setenv` | S1 — scan `libdoc/**/*.go` (skip `*_test.go`) for process Setenv/Unsetenv of SESSION/GOCACHE | **RED** (`runner.go` still Setenv) | **GREEN** |
| `functional/session-tiny-pass` | F1 — `doctest test` tiny fixture exit 0; fixture Assert checks non-empty `d.DOCTEST_SESSION_ID` | **GREEN** (today Setenv still supplies child env) | **GREEN** (via cmd.Env only) |
| `functional/cold-cache-tiny-pass` | F2 — `doctest test --cold-cache` tiny fixture exit 0; stderr announces cold-cache / GOCACHE | **GREEN** | **GREEN** (opts.GoCache → cmd.Env) |

## How to Run

```sh
doctest vet ./tests/parallel-safe/env-no-setenv/
doctest test ./tests/parallel-safe/env-no-setenv/
# Static leaf alone (fast, default suite — unlabeled):
doctest test ./tests/parallel-safe/env-no-setenv/static-source/...
# Functional locks (heavy — nested generate + go test):
doctest test --label-all ./tests/parallel-safe/env-no-setenv/functional/...
# Full tree including heavy:
doctest test --label-all ./tests/parallel-safe/env-no-setenv/
```

Classic TDD: implementer removes `os.Setenv` for session/GOCACHE in
`libdoc/runner`, holds sid/GoCache on opts, plumbs all relevant `go` tool
execs via `cmd.Env` key-replace. Do **not** change Chdir / genDir package
vars / Stdout (later phases). Existing `tests/session-inject` and
`tests/test/cold-cache` remain the deep behavior cookbook; this tree locks the
P1 anti-pattern + smoke.

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

	"github.com/xhd2015/doctest/session"
)

// Request selects static scan vs CLI subprocess.
type Request struct {
	// Op:
	//   "static_scan" — Assert scans product sources; Run is a no-op success.
	//   "cli"         — exec req.Bin with req.Args (default).
	Op string

	Args    []string
	Env     []string
	WorkDir string
	Timeout time.Duration
	Bin     string
	UseCLI  bool // true = product binary (functional e2e)

	// ModuleRoot is the doctest module root (parent of libdoc/, tests/).
	// Set by root Setup from d.DOCTEST_ROOT.
	ModuleRoot string

	// FixtureDir is the absolute path of a temp doctest tree for functional leaves.
	FixtureDir string
	// CacheHome is an isolated DOCTEST_CACHE_HOME for cold-cache leaves.
	CacheHome string
}

// Response captures CLI outcomes (static_scan leaves leave ExitCode 0).
type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	if req.Op == "static_scan" {
		return &Response{ExitCode: 0}, nil
	}

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

// moduleRootFromD resolves the repo module root from this nested tree's DOCTEST_ROOT
// (tests/parallel-safe/env-no-setenv → three levels up).
func moduleRootFromD(d *session.Doctest) string {
	return filepath.Clean(filepath.Join(d.DOCTEST_ROOT, "..", "..", ".."))
}
```
