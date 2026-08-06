# Scenario

**Feature**: `--changed` policy via core APIs (L2), help via `cli.RunWithWriter`, sparse CLI e2e (L3)

```
# L2 default: in-process filter against fixture trees
fixture tree under t.TempDir()
  -> core.DiscoverTreeCases + FilterByChangedFiles | ChangedDoctestMarkdownFiles
  -> Response{FilteredPaths, Info, MarkdownPaths}

# L2 help (help/ only): in-process CLI
cli.RunWithWriter -> doctest <subcmd> --help

# L3 e2e (not-git-repo/, label: e2e)
testbin.Ensure -> req.Bin
  -> doctest <subcmd> --changed <dir>
```

## Preconditions

- Nested root: module root is `DOCTEST_ROOT/../..` (e2e binary build only for not-git).
- L2 policy injects fixture trees + synthetic changed path lists — no product binary and no real git required.
- L2 help uses `cli.RunWithWriter` — no product binary.
- L3 e2e uses `testbin.Ensure` with the product binary (not-git-repo only).
- Leaves write fixture trees under `t.TempDir()`; never mutate the repo tree.
- Completeness: every prior `tests/changed` scenario remains as L2 or L3 smoke.

## Steps

1. Root Setup is a no-op (no binary by default).
2. L2 policy descendants create fixture trees and set `TreeDir`, `GitRoot`, `ChangedFiles`.
3. L2 help descendants set `Args` only (`<subcmd> --help`).
4. Default `Run` dispatches policy, `cli.RunWithWriter`, or binary e2e.
5. `not-git-repo/` Setup sets `UseCLI`, builds binary once via `testbin.Ensure`.

## Context

- `Request` / `Response` / `Run` are defined only in `DOCTEST.md`.
- Parallel-safe: each leaf uses `t.TempDir()`.
- **Layer**: L2 in-process policy + help is the mass; L3 e2e is sparse and labeled `e2e` when full-integration.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Default: in-process. not-git-repo/ Setup sets UseCLI + Bin.
	_ = d
	_ = req
	return nil
}
```
