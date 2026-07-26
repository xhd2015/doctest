# Go -buildvcs=auto vs local git state

## Version
0.0.1

Standalone regression tree for the Go toolchain flag `-buildvcs` (default
`auto`). Leaves build a tiny `package main` inside a **synthetic git repo**
created in a temp dir — not the doctest product binary.

## Why this tree exists

CI often uses shallow checkouts (`actions/checkout` default `fetch-depth: 1`).
That co-occurs with `error obtaining VCS status` failures, so shallow is easy
to blame. This tree locks the real axis:

| Git state | `go build -buildvcs=auto` |
|-----------|---------------------------|
| Healthy full clone | succeeds (stamps `vcs.*`) |
| Healthy shallow (`--depth 1`) | succeeds (same as full for tip) |
| Git present but status/log fails (broken HEAD) | fails with exit status 128 |

## DSN (Domain Specific Notion)

### Participants
- **`go build -buildvcs=auto`** — Go's default VCS stamping mode for main
  packages in a local module. Runs `git status --porcelain` and
  `git log -1` when it decides the main package, module, and cwd share a repo.
- **Origin module** — a minimal `example.com/buildvcsapp` repo with several
  commits, used only as a `file://` clone source.
- **Clone work tree** — full clone, shallow clone, or clone then corrupted
  `.git/HEAD`.
- **Assert** — checks exit code and known error substrings; does not invoke
  doctest product injection (doctest no longer forces `-buildvcs=false`).

### Behaviors
- **Create origin** — write `go.mod` + `main.go`, `git init`, commit twice.
- **Clone** — full or `--depth 1` into a fresh temp dir.
- **Break git** — overwrite `.git/HEAD` so `git status` exits 128.
- **Build** — `go build -buildvcs=auto -o <bin> .` in the clone.
- **Stamp check** — on success, `go version -m` shows `vcs` settings when
  git is healthy.

## Decision Tree

```
buildvcs-auto/
├── healthy-full-clone/       → full history, auto succeeds + stamps
├── healthy-shallow-clone/    → depth 1, auto succeeds + stamps (shallow ≠ fail)
└── broken-git-status/        → broken HEAD, auto fails (real fail axis)
```

## How to Run

```sh
doctest test ./tests/build/buildvcs-auto/...
```

```go
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

// Request is filled by leaf Setup: path of the clone to build.
type Request struct {
	// CloneDir is the working tree where go build runs (must contain go.mod).
	CloneDir string
	// OutBin is the -o path for go build; default CloneDir/app.bin.
	OutBin string
}

// Response captures go build outcome for Assert.
type Response struct {
	ExitCode int
	Output   string // combined stdout+stderr
	OutBin   string // binary path on success attempts
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	if req == nil || req.CloneDir == "" {
		return nil, fmt.Errorf("CloneDir is required")
	}
	outBin := req.OutBin
	if outBin == "" {
		outBin = filepath.Join(req.CloneDir, "app.bin")
	}
	cmd := exec.Command("go", "build", "-buildvcs=auto", "-o", outBin, ".")
	cmd.Dir = req.CloneDir
	// Clean env for build; inherit PATH so go/git resolve.
	cmd.Env = append(os.Environ(), "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	resp := &Response{
		Output: string(out),
		OutBin: outBin,
	}
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			resp.ExitCode = ee.ExitCode()
		} else {
			resp.ExitCode = 1
		}
		// Expected failures are asserted on ExitCode/Output; do not wrap as Run error.
		return resp, nil
	}
	resp.ExitCode = 0
	return resp, nil
}
```
