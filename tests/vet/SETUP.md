# Scenario

**Feature**: validate doctest trees via `runner.VetArgs` / `validate` / `cli.RunWithWriter` (in-process)

```
# L2: in-process vet against fixture trees
fixture tree under t.TempDir()
  -> runner.VetArgs | validate.RunWithOptions | cli.RunWithWriter
  -> Response{ExitCode, Stdout, Stderr}
```

## Preconditions

- Nested root: does not inherit workspace binary Run from `tests/DOCTEST.md`.
- All leaves inject fixture dirs / Args on `Request` — no product binary required.
- Leaves write fixture trees under `t.TempDir()`; never mutate the repo tree.
- Completeness: structure, anti-patterns, **vacuous Setup / prose-only SETUP**,
  path/argv, **layer-share** as L2 in-process.
- Help and verbose leaves are unlabeled (fast); no `testbin`.
- Layer-share fixtures use multi-leaf labeled ASSERT frontmatter (`e2e` / `heavy`).

## Steps

1. Root Setup is a no-op (no binary).
2. Descendants create fixture trees, set `Args` (and optional `WorkDir`).
3. `Run` dispatches in-process to `runner.VetArgs`, `validate.RunWithOptions`
   (verbose with injected Stdout), or `cli.RunWithWriter` (help / unknown).

## Context

- `Request` / `Response` / `Run` are defined only in `DOCTEST.md`.
- Parallel-safe: each leaf uses `t.TempDir()`; relative Args rewritten against WorkDir.
- **Layer**: L2 in-process for all leaves.
- **Implementer note (layer-share)**: skip L3 share when `opts.ChangedOnly`; constants
  `MaxL3Pct=10`, `MinLeaves=10`; L3 = ASSERT frontmatter label `e2e` only.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// In-process only: no testbin, no UseCLI binary path.
	_ = d
	_ = req
	return nil
}
```
