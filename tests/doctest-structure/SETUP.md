# Scenario

**Feature**: build the doctest binary for structure-layout integration tests

```
# shared CLI binary for structure leaves (cold-start: one go build per module)
testbin.Ensure(moduleRoot) -> $CACHE/doctest/selftest-bin/<key>/doctest

# invoke subcommands against temp trees or skill output
doctest vet|build|test|skill -> capture stdout, stderr, exit code
```

## Preconditions

- The doctest module root is two levels above this tree (`DOCTEST_ROOT/../..`).
- Each leaf builds or reuses the doctest binary via `req.Bin`.

## Steps

1. Resolve a shared doctest binary via `testbin.Ensure` (or helper `buildDoctestBin`).
2. Execute the binary with leaf-specific arguments.

## Context

- Shared helpers for writing temporary doctest trees live in root `DOCTEST.md`.
- Vet leaves create invalid or valid trees in `t.TempDir()` and run `doctest vet`.
- Integration leaves create minimal valid new-layout trees and run `doctest build` or `doctest test`.

```go
import (
	"path/filepath"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
)

func Setup(t *testing.T, req *Request) error {
	req.Timeout = 120 * time.Second
	req.Bin = testbin.Ensure(t, filepath.Join(DOCTEST_ROOT, "..", ".."))
	return nil
}
```
