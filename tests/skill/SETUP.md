# Scenario

**Feature**: skill list/show via in-process CLI (no product binary)

```
# list available skills
cli.RunWithWriter -> doctest skill --list -> skill names -> stdout buffer

# show embedded skill document
cli.RunWithWriter -> doctest skill <name> --show -> embedded .md -> stdout buffer
```

## Preconditions

- Nested root: does not inherit `tests/` binary `Run` or `testbin.Ensure`.
- All leaves are in-process via `cli.RunWithWriter` + `cli.Run`.
- No product binary build; no `label: e2e`.
- Completeness: eight leaves — list, tdd-show, tdd-lite-show, designer-show,
  implementer-show, doc-spec-show, code-spec-show, review-perf-show.

## Steps

1. Root Setup is a no-op (no binary).
2. Each leaf sets `req.Args` for the skill list or show under test.
3. `Run` calls `cli.RunWithWriter(&buf, req.Args)` and fills `Response.Stdout`.

## Context

- `Request` / `Response` / `Run` are defined only in this tree's `DOCTEST.md`.
- Parallel-safe: capture uses `withTestStdout` under `RunWithWriter` (no `os.Stdout` swap).
- **Layer**: in-process CLI for all leaves.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// In-process: no testbin, no timeout for subprocess.
	_ = d
	_ = req
	return nil
}
```
