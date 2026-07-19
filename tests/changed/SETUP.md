# Scenario

**Feature**: build the doctest binary for `--changed` integration tests

```
# build doctest from module source
go build ./cmd/doctest -> doctest binary

# invoke with --changed against ephemeral git repos
doctest <subcmd> --changed <fixture-tree> -> filter leaves by git status
```

## Preconditions

- The doctest module root is two levels above this tree (`DOCTEST_ROOT/../..`).
- Each leaf configures subcommand args and git fixture layout.

## Steps

1. Build the doctest binary from the module root.
2. Store the binary path in `req.Bin` for `Run` to execute.

## Context

- Fixture test trees are created inside ephemeral git repos by descendant `Setup` functions.
- `req.WorkDir` is set to the git repo root when running against fixtures.

```go
import (
"github.com/xhd2015/doctest/session"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Timeout = 120 * time.Second

	req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	return nil
}
```