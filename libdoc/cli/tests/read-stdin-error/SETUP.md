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
- Tests replace `os.Stdin` with a file/pipe/directory before calling `cli.Run()`.
- The current `readStdinIfPresent()` ignores errors; after the fix it propagates them.

## Steps
1. Child SETUP.md files configure `req.Args` and `req.StdinFile`.
2. Root `Run` replaces `os.Stdin` with the configured file, then calls `cli.Run(req.Args)`.
3. Stdout/stderr are captured to keep test output clean.

## Context
- These tests verify that errors from `os.Stdin.Stat()` and `io.ReadAll()` inside `readStdinIfPresent()` propagate through `cli.Run()` instead of being swallowed.

```go
import (
	"bytes"
	"io"
	"os"
	"testing"
	"github.com/xhd2015/doctest/libdoc/cli"
)

```
