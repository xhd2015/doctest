# Scenario

**Feature**: build doctest and run `skills update` in an isolated directory

```
go build ./cmd/doctest -> req.Bin
mkdir temp project dir -> optional skill install -> doctest skills update
```

## Preconditions

- Module root is `DOCTEST_ROOT/../..`.
- Implementation provides `skills` subcommand (RED until wired).

## Steps

1. Build the doctest binary into a temp path.
2. Create `req.WorkDir` as an empty project directory for subprocess cwd.
3. Run any `req.PreInstalls` CLI sequences.
4. Execute `req.Args` and capture stdout, stderr, exit code.

## Context

- Uses the same build pattern as the parent `tests/SETUP.md` tree but keeps a
  dedicated `Run` contract for skills batch update.

```go
import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
)

func Setup(t *testing.T, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 60 * time.Second
	}
	if req.Bin == "" {
		req.Bin = testbin.Ensure(t, filepath.Join(DOCTEST_ROOT, "..", ".."))
	}
	if req.WorkDir == "" {
		req.WorkDir = t.TempDir()
	}
	return nil
}
```