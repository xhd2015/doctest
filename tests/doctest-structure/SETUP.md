# Scenario

**Feature**: build the doctest binary for structure-layout integration tests

```
# build fresh doctest binary from module source
go build ./cmd/doctest -> doctest binary

# invoke subcommands against temp trees or skill output
doctest vet|build|test|skill -> capture stdout, stderr, exit code
```

## Preconditions

- The doctest module root is two levels above this tree (`DOCTEST_ROOT/../..`).
- Each leaf builds or reuses the doctest binary via `req.Bin`.

## Steps

1. Build the doctest binary from the module root.
2. Execute the binary with leaf-specific arguments.

## Context

- Shared helpers for writing temporary doctest trees live in root `DOCTEST.md`.
- Vet leaves create invalid or valid trees in `t.TempDir()` and run `doctest vet`.
- Integration leaves create minimal valid new-layout trees and run `doctest build` or `doctest test`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Timeout = 30 * time.Second
	req.Bin = buildDoctestBin(t)
	return nil
}
```