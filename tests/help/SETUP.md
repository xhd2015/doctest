# Scenario

**Feature**: top-level and scoped help via in-process CLI (no product binary)

```
# top-level usage
cli.RunWithWriter -> doctest --help -> list subcommands -> stdout buffer

# scoped help
cli.RunWithWriter -> doctest <subcmd> --help -> flags, description -> stdout buffer
```

## Preconditions

- Nested root: does not inherit `tests/` binary `Run` or `testbin.Ensure`.
- All leaves are in-process via `cli.RunWithWriter` + `cli.Run`.
- No product binary build; no `label: e2e`.
- Completeness: four leaves — top-level, test-options, build-options, agent-generate.

## Steps

1. Root Setup is a no-op (no binary).
2. Each leaf sets `req.Args` for the help variant under test.
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
