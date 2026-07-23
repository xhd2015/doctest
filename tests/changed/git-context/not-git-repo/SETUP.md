# Scenario

**Feature**: sparse L3 binary smokes — `--changed` outside a git repository

```
# build once per session
testbin.Ensure -> req.Bin
req.UseCLI = true
  -> doctest <subcmd> --changed <non-git-dir> -> hard error
```

## Preconditions

- A valid fixture test tree exists in a temp directory with **no** git repository.
- Labeled `heavy` so default discovery skips them.
- Process boundary (CLI git gate) is the SUT.

## Steps

1. Set `UseCLI` and ensure `req.Bin` via session-shared `testbin.Ensure`.
2. Create a fixture test tree without `git init`.
3. Leaf sets subcommand args with `--changed`.

## Context

- Module root: `DOCTEST_ROOT/../..`.

```go
import (
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
	"github.com/xhd2015/doctest/libdoc/testtree"
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.UseCLI = true
	if req.Timeout == 0 {
		req.Timeout = 45 * time.Second
	}
	if req.Bin == "" {
		req.Bin = testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", ".."))
	}
	dir := t.TempDir()
	testtree.WritePassFailTree(t, filepath.Join(dir, "tests"), 2, 0)
	req.WorkDir = dir
	return nil
}
```
