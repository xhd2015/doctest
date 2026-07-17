# Scenario

**Feature**: build the doctest binary and invoke it as a subprocess to test CLI behavior

```
# shared CLI binary for all leaves (cold-start: one go build per module)
testbin.Ensure(moduleRoot) -> $CACHE/doctest/selftest-bin/<key>/doctest

# invoke binary as subprocess, capture everything
doctest <args> -> {stdout, stderr, exit code}
```

## Preconditions
- The doctest module root is the parent of this test tree (`DOCTEST_ROOT/..`).
- The tests are executed by the doc-style test runner from this test tree.
- Each leaf sets the doctest arguments it wants to execute.

## Steps
1. Resolve a shared doctest binary via `testbin.Ensure` (build once; reuse across leaves).
2. Execute the binary given by `req.Bin`.
3. Capture stdout, stderr, exit code, and the raw execution error.

## Context
- These are real integration tests, not mocked unit tests.
- Agent tests expect `fake-codex` to be in PATH.

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
	"github.com/xhd2015/doctest/libdoc/testbin"
)

func Setup(t *testing.T, req *Request) error {
	req.Timeout = 30 * time.Second

	req.Bin = testbin.Ensure(t, filepath.Join(DOCTEST_ROOT, ".."))
	return nil
}
```
