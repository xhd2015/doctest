# Scenario

**Feature**: stdin read errors propagate through agent design and implement commands

```
# optional stdin before agent dispatch
cli.Run(args) -> readStdinIfPresent -> runAgentDesign | runAgentImplement

# errors must not be swallowed
broken stdin -> Stat/ReadAll error -> returned to caller
```

## Preconditions
- The `cli` package is importable.
- Tests inject stdin via `cli.TestExported_RunWithStdin` (never reassign process stdin).
- Product `readStdinIfPresent` honors the inject and propagates Stat/ReadAll errors.

## Steps
1. Child SETUP.md files configure `req.Args` and `req.StdinFile`.
2. Root `Run` calls `cli.TestExported_RunWithStdin(req.Args, req.StdinFile)`.
3. No process stdio swap — Parallel-safe harness.

## Context
- These tests verify that errors from stdin `Stat()` and `io.ReadAll()` inside `readStdinIfPresent()` propagate through `cli.Run()` instead of being swallowed.

```go
import (
	"testing"
)
```
