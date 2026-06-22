# Scenario

**Feature**: the doctest module root is two levels above this test tree (`DOCTEST_ROOT/../..`)

```
# logf formats agent session events for display
doctest agent logf <session-id> -> reads event file -> formatted text -> stdout

# show-status reports session progress
doctest agent show-status <session-id> -> session state -> stdout
```

## Preconditions
- The doctest module root is two levels above this test tree (`DOCTEST_ROOT/../..`).
- Tests verify that `traceSession` and `showStatus` output lines use timestamped `Logf` while UI framing uses bare `fmt.Fprintf`.

## Steps
1. Build the doctest binary.
2. Shell out to the doctest binary with `req.Args`, capturing stdout, stderr, and exit code.
3. Leaves configure inputs via `req.Args` or `req.Env`.

## Context
- These tests verify the contract: Logf produces `[2006-01-02T15:04:05]` prefixed output.
- `traceSession` and `showStatus` tests create real session directories and run the doctest binary.
- The `logf/` subtree has its own `DOCTEST.md` root because it calls `subagent.Logf` in-process.

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
	libdocbuild "github.com/xhd2015/doctest/libdoc/build"
)

func Setup(t *testing.T, req *Request) error {
	req.Timeout = 30 * time.Second

	tmp := t.TempDir()
	doctestBin := filepath.Join(tmp, "doctest")
	buildDir := filepath.Join(DOCTEST_ROOT, "..", "..")
	buildArgs := []string{"build", "-o", doctestBin}
	if libdocbuild.NeedsBuildVCSFlag(buildDir) {
		buildArgs = append(buildArgs, "-buildvcs=false")
	}
	buildArgs = append(buildArgs, "./cmd/doctest")
	build := exec.Command("go", buildArgs...)
	build.Dir = buildDir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build doctest: %v\n%s", err, string(out))
	}
	req.Bin = doctestBin
	return nil
}
```
