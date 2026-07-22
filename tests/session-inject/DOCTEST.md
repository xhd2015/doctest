# Embedded Session Local Module — Integration Tests

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Doctest CLI** — discovers doctest leaves, assembles generated Go, detects
  `github.com/xhd2015/doctest/session` imports, and chooses compile strategy.
- **sessionmod** — embedded session package sources shipped inside the doctest
  binary (`libdoc/sessionmod`).
- **Session cache** — content-addressed directory under
  `$UserCacheDir/doctest/session-mod/<md5>/` holding session sources + standalone
  `go.mod` with `module github.com/xhd2015/doctest/session`.
- **Nested testcase module** — generated `module testcase` `go.mod` may append
  `replace github.com/xhd2015/doctest/session => <cache>` when a leaf imports session.
- **Consumer leaf** — SETUP/ASSERT that imports `session` and calls `Once`.

**Behaviors**

- Session package is always required for inject (`d *session.Doctest`) →
  materialize session-mod cache (write-once) before code generation (even when
  author SETUP/ASSERT do not import session).
- Nested module + session (always) → `replace` for session points at cache.
- Import path stays `github.com/xhd2015/doctest/session` (no rewrite).
- A leaf that imports session and calls `Once` can compile and run when the
  embedded module is replaced correctly and `DOCTEST_SESSION_ID` is set for the
  subprocess.

## Decision Tree

```
session-inject/
├── cache/                                    [session-mod materialization]
│   ├── first-run-materializes/               B1: creates session-mod/<md5>/
│   ├── second-run-idempotent/                B2: second run does not rewrite
│   └── no-import-skips/                      B3: no author session import → still materialize (inject)
└── replace/                                  [go.mod replace + consumer Once]
    ├── replace-in-gomod/                     R1: replace session => cache path
    └── once-call-succeeds/                   R2: leaf imports session; Once works
```

## Test Index

| Leaf | Scenario |
|------|----------|
| `cache/first-run-materializes` | B1 — first run creates `$CACHE/doctest/session-mod/<md5>/` with go.mod + sources |
| `cache/second-run-idempotent` | B2 — second run leaves cache bytes unchanged |
| `cache/no-import-skips` | B3 — run without author session import still creates session-mod (inject) |
| `replace/replace-in-gomod` | R1 — nested go.mod contains session replace pointing at cache |
| `replace/once-call-succeeds` | R2 — subprocess leaf imports session and Once returns JSON |

## How to Run

```sh
doctest vet ./tests/session-inject/
doctest test ./tests/session-inject/          # expect RED until sessionmod + materialize land
doctest test ./tests/session-inject/cache/...
doctest test ./tests/session-inject/replace/...
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
)

type Request struct {
	Args    []string
	Env     []string
	WorkDir string
	Timeout time.Duration
	Bin     string
	// GenDir is the absolute --gen-dir path for this leaf (request-local; no package var).
	GenDir  string
	// ModuleRoot / TestDir are the temp fixture paths for this leaf (request-local).
	ModuleRoot string
	TestDir    string
	// GoModBefore is a snapshot of session-mod go.mod before a measured run (request-local).
	GoModBefore []byte
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
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
